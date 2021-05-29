package app

import (
	"github.com/gin-gonic/gin"

	"github.com/mendezdev/tgo_flattener/ping"
)

func routes(h handlers) *gin.Engine {
	router := gin.Default()

	router.GET("/ping", ping.Ping)
	router.POST("/flat", h.Flat.Post)

	return router
}
