package api

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	ginlogrus "github.com/toorop/gin-logrus"

	"word-search/pkg/logger"
)

// InitServer initialize server
func InitServer() {
	const section = "server.InitServer"
	logger.Log.Infoln(section, "starting")

	router := gin.Default()
	router.SetTrustedProxies(nil)

	router.Use(ginlogrus.Logger(logger.Log), gin.Recovery())

	Router(router)

	err := router.Run(fmt.Sprintf(":%s", os.Getenv("SERVER_PORT")))
	if err != nil {
		logger.Log.Errorln(section, err)

		return
	}

	logger.Log.Infoln(section, "finished")
}
