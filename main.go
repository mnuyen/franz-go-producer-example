package main

import (
	"fmt"
	"net/http"
	"os"

	"com.mmuyen.go.franz.producer/app"
	"github.com/gin-gonic/gin"
)

func main() {
	app.StartProducer()
	go func() {
		Controller()
	}()
	select {}
}

func Controller() {
	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	router.POST("/send/:topic", func(context *gin.Context) {
		topic := context.Param("topic")
		jsonData, err := context.GetRawData()
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		app.SendKafkaMessage(topic, "test", string(jsonData))
		context.JSON(200, gin.H{
			"message": "success",
		})
	})
	port := os.Getenv("port")
	if port == "" {
		port = ":8091"
	}
	serverError := http.ListenAndServe(port, router)
	if serverError != nil {
		fmt.Println(serverError.Error())
	}
}
