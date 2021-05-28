package flattener

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler interface {
	Post(c *gin.Context)
}

type handler struct {
	gtw Gateway
}

func NewHandler(flatGateway Gateway) Handler {
	return &handler{
		gtw: flatGateway,
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
