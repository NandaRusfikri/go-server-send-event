package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type Request struct {
	Message string `json:"message"`
}
type DataEvent struct {
	Message string `json:"message"`
	Time    string `json:"time"`
}

var Event = make(chan DataEvent)

func main() {
	router := gin.Default()

	router.Static("/client", "./client")
	router.GET("/api/polling", HandlePollng)
	router.GET("/api/sse", HandleSSE)
	router.POST("/api/sse", SendData)
	router.Run(":4444")
}
func HandlePollng(c *gin.Context) {
	rand.NewSource(time.Now().UnixNano())
	randomNumber := rand.Intn(10)
	c.JSON(http.StatusOK, gin.H{
		"time":    time.Now().Format("15:04:05"),
		"message": randomNumber,
	})
}

func HandleSSE(c *gin.Context) {
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			rand.NewSource(time.Now().UnixNano())
			randomNumber := rand.Intn(10)
			jsonData, err := json.Marshal(DataEvent{
				Message: strconv.Itoa(randomNumber),
				Time:    time.Now().Format("15:04:05"),
			})
			if err != nil {
				fmt.Println("Error encoding JSON:", err)
				return
			}
			c.Writer.WriteString(fmt.Sprintf("event: %s\ndata: %s\n\n", "event-update", jsonData))
			c.Writer.Flush()
		case <-c.Request.Context().Done():
			fmt.Println("Client closed the connection.")
			return
		}
	}
}

//func HandleSSE(c *gin.Context) {
//	c.Header("Content-Type", "text/event-stream")
//	c.Header("Cache-Control", "no-cache")
//	c.Header("Connection", "keep-alive")
//
//	for {
//		select {
//		case data := <-Event:
//			jsonData, err := json.Marshal(data)
//			if err != nil {
//				fmt.Println("Error encoding JSON:", err)
//				return
//			}
//			c.Writer.WriteString(fmt.Sprintf("event: %s\ndata: %s\n\n", "event-update", jsonData))
//			c.Writer.Flush()
//		case <-c.Request.Context().Done():
//			fmt.Println("Client closed the connection.")
//			return
//
//		}
//	}
//}

func SendData(c *gin.Context) {
	var input Request

	if err := c.ShouldBindJSON(&input); err != nil {
		c.AbortWithStatusJSON(400, err.Error())
		return
	}

	ServerTime := time.Now().Format("15:04:05")
	Event <- DataEvent{
		Message: input.Message,
		Time:    ServerTime,
	}

	c.JSON(200, map[string]interface{}{
		"message": "oke",
		"time":    ServerTime,
	})

}
