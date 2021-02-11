package httpServer

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
)

const (
	AdminLogin = "admin"
)

func InitServer() error {
	log.Println("Initializing http-server STARTED")

	router := gin.Default()

	randomUuid, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	adminPassword := randomUuid.String()
	log.Println("Admin password: ", adminPassword)

	authorized := router.Group("/", gin.BasicAuth(gin.Accounts{
		AdminLogin: adminPassword,
	}))

	obfuscatorRouter(*authorized)

	err = router.Run()
	if err != nil {
		return err
	}

	log.Println("Initializing http-server FINISHED")
	return nil
}
