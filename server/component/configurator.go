package component

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/x-research-team/bus"
	"github.com/x-research-team/contract"
)

func Configure() contract.ComponentModule {
	return func(c contract.IComponent) {
		component := c.(*Component)
		component.engine = gin.Default()
		component.engine.POST("/api", func(ctx *gin.Context) {
			buffer, err := ioutil.ReadAll(ctx.Request.Body)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err})
				return
			}
			m := new(KernelMessage)
			if err = json.Unmarshal(buffer, m); err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"Error": err})
				return
			}
			message := bus.Message(m.Route, m.Command, m.Message)
			component.trunk <- bus.Signal(message)
			ctx.JSON(http.StatusOK, gin.H{"ID": message.ID()})
		})
		c = component
	}
}

type KernelMessage struct{ Route, Command, Message string }
