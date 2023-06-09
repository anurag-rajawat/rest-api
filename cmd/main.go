package main

import (
	"runtime"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"github.com/anurag-rajawat/rest-api/pkg/routes"
	"github.com/anurag-rajawat/rest-api/pkg/utils"
)

func init() {
	printVersion()
	utils.Init()
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	routes.RegisterRoutes(router)

	log.Info("Listening and serving HTTP on :8080")
	err := router.Run()
	if err != nil {
		log.Fatalf("Failed to start server due to %s", err)
	}
}

func printVersion() {
	log.Infof("go version: %s", runtime.Version())
	log.Infof("go os/arch: %s/%s", runtime.GOOS, runtime.GOARCH)
}
