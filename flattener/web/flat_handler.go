package flattener

import (
	"net/http"

	"github.com/gin-gonic/gin"
	flatGtw "github.com/mendezdev/tgo_flattener/flattener/gateway"
)

type FlatHandler interface {
	Flat(c *gin.Context)
}

type flatHandlerImpl struct {
	gtw flatGtw.FlatService
}

func NewFlatHandler() FlatHandler {
	return &flatHandlerImpl{
		gtw: flatGtw.NewFlatService(),
	}
}

func (h *flatHandlerImpl) Flat(c *gin.Context) {
	var unflatted []interface{}
	if err := c.ShouldBindJSON(&unflatted); err != nil {
		c.JSON(http.StatusBadRequest, "error parsing body")
		return
	}

	flatResponse := h.gtw.FlatResponse(unflatted)
	c.JSON(http.StatusOK, flatResponse)
}
