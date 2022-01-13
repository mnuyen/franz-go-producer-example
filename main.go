package main

import (
	"com.mmuyen.go.franz.producer/docs"
	"fmt"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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
	docs.SwaggerInfo.BasePath = "/"
	router.GET("/ping", pingController)
	router.POST("/send/:topic", sendElementToKafkaBroker)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	port := os.Getenv("port")
	if port == "" {
		port = ":8091"
	}

	serverError := http.ListenAndServe(port, router)
	if serverError != nil {
		fmt.Println(serverError.Error())
	}
}

// @BasePath /ping

// PingExample godoc
// @Summary ping example
// @Schemes
// @Description do ping
// @Tags example
// @Accept json
// @Produce json
// @Success 200 {string} Helloworld
// @Router /ping [get]
func pingController(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

// @BasePath /send/:topic

// Send Topic Example godoc
// @Summary ping example
// @Schemes
// @Description do ping
// @Tags example
// @Accept json
// @Produce json
// @Success 200 {string} Helloworld
// @Router /send/:topic [post]
func sendElementToKafkaBroker(context *gin.Context) {
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
}
