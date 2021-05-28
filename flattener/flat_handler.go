package flattener

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mendezdev/tgo_flattener/apierrors"
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
		apiErr := apierrors.NewBadRequestError("error parsing body")
		c.JSON(http.StatusBadRequest, apiErr)
		return
	}

	flatResponse := h.gtw.FlatResponse(unflatted)
	c.JSON(http.StatusOK, flatResponse)
}
