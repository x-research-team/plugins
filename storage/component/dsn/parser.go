package dsn

import (
	"github.com/x-research-team/bus"
	"github.com/x-research-team/utils/file"
)

func Parse() map[string]string {
	v := make(map[string]string)
	if err := file.Read("config", "storage.json", &v); err != nil {
		bus.Error <- err
		return v
	}
	return v
}
