package handlers_test

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/anurag-rajawat/rest-api/pkg/types"
)

var (
	Db      *gorm.DB
	cleanup func()
)

func TestMain(m *testing.M) {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatalf("Could not load env variables: %s", err)
	}
	Db, cleanup = setupDB()
	code := m.Run()
	// You can't defer this because os.Exit doesn't care for defer
	cleanup()
	os.Exit(code)
}

// GetTestGinContext returns a gin context for testing
func GetTestGinContext(w *httptest.ResponseRecorder) *gin.Context {
	gin.SetMode(gin.TestMode)
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = &http.Request{
		Header: make(http.Header),
		URL:    &url.URL{},
	}
	return ctx
}

// setupDB return a fresh postgres database for testing
func setupDB() (*gorm.DB, func()) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not construct pool: %s", err)
	}

	err = pool.Client.Ping()
	if err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
	}

	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "15",
		Env:        []string{"POSTGRES_PASSWORD=" + os.Getenv("DB_PASSWORD"), "POSTGRES_DB=" + os.Getenv("DB_NAME")},
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.NeverRestart()
	})

	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	dsn := fmt.Sprintf("host=localhost user=postgres password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		resource.GetPort("5432/tcp"),
	)

	log.Info("Connecting to database")

	var gdb *gorm.DB
	err = pool.Retry(func() error {
		gdb, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		})
		if err != nil {
			return err
		}
		db, err := gdb.DB()
		if err != nil {
			return err
		}
		return db.Ping()
	})
	log.Info("Successfully connected to database")

	return gdb, func() {
		if err := pool.Purge(resource); err != nil {
			log.Fatalf("Could not purge resource: %s", err)
		}
	}
}

func InitDb() error {
	err := Db.AutoMigrate(&types.User{})
	if err != nil {
		return errors.New(fmt.Sprintf("Could not create relation: %s", err))
	}
	return nil
}

func SeedOneUser() types.User {
	err := InitDb()
	if err != nil {
		log.Fatal(err)
	}

	user := types.User{
		UserName: "testuser",
		Email:    "test@gmail.com",
		Password: "passwd",
	}
	user.Password = hashedPasswd(user.Password)

	err = Db.Model(&types.User{}).Create(&user).Error
	if err != nil {
		log.Fatalf("Could not seed user: %s", err)
	}
	return user
}

func CleanDb() {
	err := Db.Exec(`DROP SCHEMA public CASCADE;CREATE SCHEMA public;`).Error
	if err != nil {
		log.Fatalf("Could not clean database: %s", err)
	}
}

func hashedPasswd(password string) string {
	hashedPasswd, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPasswd)
}
