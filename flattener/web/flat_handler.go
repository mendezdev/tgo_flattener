package flattener

import (
	"net/http"

	"github.com/gin-gonic/gin"
	flatGtw "github.com/mendezdev/tgo_flattener/flattener/gateway"
)

type Handler interface {
	Post(c *gin.Context)
}

type handler struct {
	gtw flatGtw.Gateway
}

func NewHandler() Handler {
	return &handler{
		gtw: flatGtw.NewGateway(),
	}
}

func (h *handler) Post(c *gin.Context) {
	var unflatted []interface{}
	if err := c.ShouldBindJSON(&unflatted); err != nil {
		c.JSON(http.StatusBadRequest, "error parsing body")
		return
	}

	flatResponse := h.gtw.FlatResponse(unflatted)
	c.JSON(http.StatusOK, flatResponse)
}
