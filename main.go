package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

//go:embed index.html
var indexHTML []byte

type Request struct {
	Message string `json:"message"`
}
type DataEvent struct {
	ServerTime string `json:"server_time"`
	Message    string `json:"message"`
}

var Event = make(chan DataEvent)

func main() {
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html; charset=utf-8", indexHTML)
	})

	router.GET("/events", HandleEvent)
	router.POST("/send_data", SendData)
	router.Run(":4444")
}

func SendData(c *gin.Context) {
	var input Request

	if err := c.ShouldBindJSON(&input); err != nil {
		c.AbortWithStatusJSON(400, err.Error())
		return
	}

	ServerTime := time.Now().Format("2006-02-01 15:04")
	Event <- DataEvent{
		Message:    input.Message,
		ServerTime: ServerTime,
	}

	c.JSON(200, map[string]interface{}{
		"status": "oke",
		"date":   ServerTime,
	})

}

func HandleEvent(c *gin.Context) {
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	for {
		select {
		case data := <-Event:
			fmt.Println("data ", data)
			jsonData, err := json.Marshal(data)
			if err != nil {
				fmt.Println("Error encoding JSON:", err)
				return
			}
			message := fmt.Sprintf("event: %s\ndata: %s\n\n", "event-update", jsonData)

			c.Writer.WriteString(message)
			c.Writer.Flush()
		case <-c.Request.Context().Done():
			fmt.Println("Client closed the connection.")
			return

		}
	}
}

//func HandleEvent(c *gin.Context) {
//	c.Header("Content-Type", "text/event-stream")
//	c.Header("Cache-Control", "no-cache")
//	c.Header("Connection", "keep-alive")
//
//	ticker := time.NewTicker(3 * time.Second)
//	defer ticker.Stop()
//
//	for {
//		select {
//		case <-ticker.C:
//			eventType := "price-update"
//			rand.Seed(time.Now().UnixNano())
//			randomNumber := rand.Intn(10)
//
//			data := map[string]interface{}{
//				"message":      randomNumber,
//				"ServerTime": time.Now().Format("15:04:05"),
//			}
//			jsonData, err := json.Marshal(data)
//			if err != nil {
//				fmt.Println("Error encoding JSON:", err)
//				return
//			}
//			message := fmt.Sprintf("event: %s\ndata: %s\n\n", eventType, jsonData)
//
//			c.Writer.WriteString(message)
//			c.Writer.Flush()
//		case <-c.Request.Context().Done():
//			fmt.Println("Client closed the connection.")
//			return
//		}
//	}
//}
