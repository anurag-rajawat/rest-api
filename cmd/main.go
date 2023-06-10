package main

import (
	"os"
	"runtime"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"

	"github.com/anurag-rajawat/rest-api/pkg/database"
	"github.com/anurag-rajawat/rest-api/pkg/routes"
)

type config struct {
	dbHost     string
	dbUser     string
	dbPassword string
	dbName     string
	dbPort     string
}

func init() {
	printVersion()
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Failed to load env variables due to: %s", err)
	}
}

func main() {
	cfg := initConfig()
	db, err := database.ConnectToDb(cfg.dbHost, cfg.dbUser, cfg.dbPassword, cfg.dbName, cfg.dbPort)
	if err != nil {
		log.Fatalf("Failed to connect to database due to: %s", err)
	}

	if err := database.Initialize(db); err != nil {
		log.Fatal(err)
	}

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	routes.RegisterRoutes(router, db)

	log.Info("Listening and serving HTTP on :8080")
	err = router.Run()
	if err != nil {
		log.Fatalf("Failed to start server due to %s", err)
	}
}

func initConfig() config {
	dbHost, ok := os.LookupEnv("DB_HOST")
	if !ok {
		log.Fatal("DB host is not specified")
	}

	dbUser, ok := os.LookupEnv("DB_USER")
	if !ok {
		log.Fatal("DB User is not specified")
	}

	dbPasswd, ok := os.LookupEnv("DB_PASSWORD")
	if !ok {
		log.Fatal("DB Password is not specified")
	}

	dbName, ok := os.LookupEnv("DB_NAME")
	if !ok {
		log.Fatal("DB Name is not specified")
	}

	dbPort, ok := os.LookupEnv("DB_PORT")
	if !ok {
		log.Warn("DB Port is not specified. Defaults to 5432")
		dbPort = "5432"
	}
	return config{
		dbHost:     dbHost,
		dbUser:     dbUser,
		dbPassword: dbPasswd,
		dbName:     dbName,
		dbPort:     dbPort,
	}
}

func printVersion() {
	log.Infof("go version: %s", runtime.Version())
	log.Infof("go os/arch: %s/%s", runtime.GOOS, runtime.GOARCH)
}
