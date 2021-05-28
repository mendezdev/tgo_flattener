package app

import (
	"github.com/gin-gonic/gin"

	flattener "github.com/mendezdev/tgo_flattener/flattener/web"
)

func routes() *gin.Engine {
	router := gin.Default()

	flatHandler := flattener.NewFlatHandler()
	router.POST("/flat", flatHandler.Flat)

	return router
}
