package app

import (
	"github.com/gin-gonic/gin"

	flattener "github.com/mendezdev/tgo_flattener/flattener/web"
)

func routes() *gin.Engine {
	router := gin.Default()

	flatService := flattener.NewFlatService()
	router.POST("/flat", flatService.FlatHandler)

	return router
}
