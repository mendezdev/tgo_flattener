package app

import (
	"github.com/gin-gonic/gin"

	flattener "github.com/mendezdev/tgo_flattener/flattener/web"
	"github.com/mendezdev/tgo_flattener/ping"
)

func routes() *gin.Engine {
	router := gin.Default()

	router.GET("/ping", ping.Ping)

	flatHandler := flattener.NewFlatHandler()
	router.POST("/flat", flatHandler.Flat)

	return router
}
