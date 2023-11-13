package main

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

//go:embed index.html
var indexHTML []byte

func main() {

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
				event, err := formatServerSentEvent("price-update", time.Now().Format("15:04:05"))
				if err != nil {
					fmt.Println(err)

				}
				fmt.Println("event ", event)
				fmt.Fprintf(w, event)
				w.(http.Flusher).Flush()
			case <-r.Context().Done():
				fmt.Println("Client closed the connection.")
				return
			}
		}

	})

	http.ListenAndServe(":4444", nil)
}

func formatServerSentEvent(event string, data any) (string, error) {
	m := map[string]any{
		"data": data,
	}

	buff := bytes.NewBuffer([]byte{})

	encoder := json.NewEncoder(buff)

	err := encoder.Encode(m)
	if err != nil {
		return "", err
	}

	sb := strings.Builder{}

	sb.WriteString(fmt.Sprintf("event: %s\n", event))
	sb.WriteString(fmt.Sprintf("data: %v\n\n", buff.String()))

	return sb.String(), nil
}
