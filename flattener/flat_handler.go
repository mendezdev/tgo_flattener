package flattener

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mendezdev/tgo_flattener/apierrors"
)

type Handler interface {
	Post(c *gin.Context)
	GetAll(c *gin.Context)
}

type handler struct {
	gtw Gateway
}

func NewHandler(flatGateway Gateway) Handler {
	return &handler{
		gtw: flatGateway,
	}
}

// Post will flat the request array
// only is available to receive arrays of simple mixed values
func (h *handler) Post(c *gin.Context) {
	var unflatted []interface{}
	if err := c.ShouldBindJSON(&unflatted); err != nil {
		apiErr := apierrors.NewBadRequestError("error parsing body")
		c.JSON(http.StatusBadRequest, apiErr)
		return
	}

	flatResponse, err := h.gtw.FlatResponse(unflatted)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}
	c.JSON(http.StatusOK, flatResponse)
}

// GetAll it will return a FlatInfo with a limit.
// You can see the limit configured in config/config.go file
func (h *handler) GetAll(c *gin.Context) {
	flats, err := h.gtw.GetFlats()
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusOK, flats)
}
