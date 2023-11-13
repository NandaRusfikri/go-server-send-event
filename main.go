package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

//go:embed index.html
var indexHTML []byte

func main() {
	Images := []string{
		"https://www.banksinarmas.com/id/public/revamp/logoj.png",
		"https://bengkuluekspress.disway.id/upload/2014/01/bank-sinarmas-syariah.jpg",
		"https://statik.tempo.co/data/2010/12/20/id_57840/57840_620.jpg",
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(indexHTML)
	})

	http.HandleFunc("/crypto-price", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		ticker := time.NewTicker(3 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:

				eventType := "price-update"
				rand.Seed(time.Now().UnixNano())

				randomNumber := rand.Intn(3)

				data := map[string]interface{}{
					"image":      Images[randomNumber],
					"ServerTime": time.Now().Format("15:04:05"),
				}
				jsonData, err := json.Marshal(data)
				if err != nil {
					fmt.Println("Error encoding JSON:", err)
					return
				}
				message := fmt.Sprintf("event: %s\ndata: %s\n\n", eventType, jsonData)

				fmt.Fprintf(w, message)

				w.(http.Flusher).Flush()
			case <-r.Context().Done():
				fmt.Println("Client closed the connection.")
				return
			}
		}

	})

	http.ListenAndServe(":4444", nil)
}
