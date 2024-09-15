package main

import (
	"natsdelivery/natsdelivery"

	"github.com/gin-gonic/gin"
	conf "github.com/kowiste/boilerplate/src/config"
	"github.com/kowiste/boilerplate/src/messaging/nats"

	"github.com/kowiste/config"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// InitializeDB initializes the database connection
func InitializeDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

var nd = &natsdelivery.NatsDelivery{}

func main() {
	// Initialize the database connection
	db, err := InitializeDB()
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}
	err = config.New[conf.BoilerConfig](config.GetPathEnv())
	if err != nil {
		panic("failed to get config: " + err.Error())
	}
	_, err = config.Get[conf.BoilerConfig]()
	if err != nil {
		panic("failed to get config: " + err.Error())
	}
	//Init messaging
	n := nats.New()
	err = n.Init()
	if err != nil {
		panic("failed to connect nats: " + err.Error())
	}
	nd = natsdelivery.New(db, n)
	err = nd.Init()
	if err != nil {
		panic("failed to init nats: " + err.Error())
	}
	// Set up Gin router
	router := gin.Default()

	// Define the POST endpoint to send a message
	router.POST("/send", sendData)

	// Start the server
	router.Run(":8080")

}
