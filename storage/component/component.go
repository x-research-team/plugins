package component

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/x-research-team/bus"
	"github.com/x-research-team/contract"
)

const (
	name  = "Storage"
	route = "storage"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

// Component
type Component struct {
	bus chan []byte

	components map[string]contract.IComponent
	trunk      contract.ISignalBus
	route      string
	uuid       string

	client map[string]*sql.DB
	fails  []error
}

// New Создать экземпляр компонента сервиса биллинга
func New(opts ...contract.ComponentModule) contract.KernelModule {
	component := &Component{
		bus:        make(chan []byte),
		components: make(map[string]contract.IComponent),
		route:      route,
		trunk:      make(contract.ISignalBus),
		client:     make(map[string]*sql.DB),
	}
	for _, o := range opts {
		o(component)
	}
	if len(component.fails) > 0 {
		for _, err := range component.fails {
			bus.Error <- fmt.Errorf("[%s] %v", name, err)
		}
		return func(service contract.IService) {
		}
	}
	bus.Add(component.trunk)
	bus.Info <- fmt.Sprintf("[%v] Initialized", name)
	return func(c contract.IService) {
		c.AddComponent(component)
		bus.Info <- fmt.Sprintf("[%v] attached to Billing Service", name)
	}
}

func (component *Component) AddComponent(c contract.IComponent) {
	component.components[c.Name()] = c
}

// Send Отправить сигнал в ядро
func (component *Component) Send(message contract.IMessage) {
	component.trunk.Send(bus.Signal(message))
}

// AddPlugin Добавить плагин на горячем ходу
func (component *Component) AddPlugin(p, name string) error {
	return nil
}

// RemovePlugin Удалить плагин на горячем ходу
func (component *Component) RemovePlugin(name string) error {
	return nil
}

// Configure Конфигурация компонета платежной системы
func (component *Component) Configure() error {
	bus.Info <- fmt.Sprintf("[%v] is configured", name)
	return nil
}

// Run Запуск компонента платежной системы
func (component *Component) Run() error {
	bus.Info <- fmt.Sprintf("[%v] component started", name)

	component.uuid = uuid.New().String()

	for {
		select {
		case data := <-component.bus:
			fmt.Printf("%s\n", data)
			m := new(KernelMessage)
			if err := json.Unmarshal(data, &m); err != nil {
				bus.Error <- err
				continue
			}
			command := new(TCommand)
			if err := json.Unmarshal(m.Data, &command); err != nil {
				bus.Error <- err
				continue
			}
			if command.Service == "" {
				bus.Error <- fmt.Errorf("unknown service")
				continue
			}
			if command.SQL == "" {
				bus.Error <- fmt.Errorf("missing sql raw")
				continue
			}
			tx, err := component.client[command.Service].Begin()
			if err != nil {
				bus.Error <- err
				continue
			}
			stmt, err := tx.Prepare(command.SQL)
			if err != nil {
				component.rollback(tx, err)
				continue
			}
			if strings.HasPrefix(strings.ToLower(command.SQL), "select") {
				rows, err := stmt.Query()
				if err != nil {
					component.rollback(tx, err)
					continue
				}
				v := make([]map[string]interface{}, 0)
				if err := rows.Scan(&v); err != nil {
					component.rollback(tx, err)
					continue
				}
				buffer, err := json.Marshal(v)
				if err != nil {
					component.rollback(tx, err)
					continue
				}
				fmt.Println(string(buffer))
			} else {
				_, err := stmt.Exec()
				if err != nil {
					component.rollback(tx, err)
					continue
				}
			}
			if err := tx.Commit(); err != nil {
				bus.Error <- err
				continue
			}
		default:
			continue
		}
	}
}

func (component *Component) Route() string { return component.route }

type KernelMessage struct {
	ID   uuid.UUID
	Data []byte
}

func (component *Component) Write(message contract.IMessage) error {
	if message.Route() != component.Route() {
		return nil
	}
	bus.Debug <- fmt.Sprintf("%#v", message)
	buffer, err := json.Marshal(&KernelMessage{
		ID:   message.ID(),
		Data: []byte(message.Data()),
	})
	if err != nil {
		return err
	}
	component.bus <- buffer
	return nil
}

func (component *Component) Read() string {
	return ""
}

func (component *Component) Pid() string {
	return component.uuid
}

func (component *Component) Name() string {
	return name
}

func (component *Component) Up(graceful bool) error {
	return nil
}

func (component *Component) Down(graceful bool) error {
	return nil
}

func (component *Component) Sleep(time.Duration) error {
	return nil
}

func (component *Component) Restart(graceful bool) error {
	return nil
}

func (component *Component) Pause() error {
	return nil
}

func (component *Component) Cron(rule string) error {
	return nil
}

func (component *Component) Stop() error {
	return nil
}

func (component *Component) Kill() error {
	return nil
}

func (component *Component) Sync(with string) error {
	return nil
}

func (component *Component) Backup(to string) error {
	return nil
}

func (component *Component) rollback(tx *sql.Tx, err error) {
	if err := tx.Rollback(); err != nil {
		bus.Error <- err
		return
	}
	bus.Error <- err
}
