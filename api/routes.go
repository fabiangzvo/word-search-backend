package api

import (
	"github.com/gin-gonic/gin"

	"word-search/pkg/logger"
	"word-search/sockets"
)

func healthCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "ok",
	})
}

func noRoute(c *gin.Context) {
	c.JSON(404, gin.H{"message": "Service not found"})
}

func handleSocket(hub *sockets.Hub) func(*gin.Context) {
	return func(c *gin.Context) {
		gameID := c.Param("gameId")

		sockets.ServeWS(c, gameID, hub)
	}
}

// Router load all available routes
func Router(router *gin.Engine) {
	const section = "api.Router"
	logger.Log.Infoln(section, "starting")

	hub := sockets.NewHub()
	go hub.Run()

	router.GET("/ws/:gameId", handleSocket(hub))
	router.GET("/health-check", healthCheck)
	router.NoRoute(noRoute)

	logger.Log.Infoln(section, "finished")
}
