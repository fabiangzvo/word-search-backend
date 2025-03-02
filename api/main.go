package api

import (
	"fmt"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	ginlogrus "github.com/toorop/gin-logrus"

	"word-search/pkg/logger"
	"word-search/sockets"
)

// InitServer initialize server
func InitServer() {
	const section = "server.InitServer"
	logger.Log.Infoln(section, "starting")

	router := gin.Default()
	router.SetTrustedProxies(nil)

	router.Use(ginlogrus.Logger(logger.Log), gin.Recovery())
	// Enable CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowCredentials: true,
		AllowHeaders:     []string{"Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "accept", "origin", "Cache-Control", "X-Requested-With", "API_KEY"},
		AllowMethods:     []string{"PUT", "PATCH", "POST", "OPTIONS", "GET", "DELETE"},
	}))

	sockets.ServeWS(router)
	Router(router)

	err := router.Run(fmt.Sprintf(":%s", os.Getenv("SERVER_PORT")))
	if err != nil {
		logger.Log.Errorln(section, err)

		return
	}

	logger.Log.Infoln(section, "finished")
}
