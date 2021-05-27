package flattener

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type FlatService interface {
	FlatHandler(c *gin.Context)
}

type flatServiceImpl struct{}

func NewFlatService() FlatService {
	return &flatServiceImpl{}
}

func (s *flatServiceImpl) FlatHandler(c *gin.Context) {
	c.JSON(http.StatusOK, "HELLO!")
}
