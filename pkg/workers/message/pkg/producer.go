package pkg

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
)

type Producer interface {
	Start() error
}

type HTTPProducer struct {
	listenAddr string
	server     *Server
	producech  chan<- *Message
}

func NewHTTPProducer(listenAddr string, producech chan *Message) *HTTPProducer {
	return &HTTPProducer{
		listenAddr: listenAddr,
		producech:  producech,
	}
}

func (p *HTTPProducer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		path  = strings.TrimPrefix(r.URL.Path, "/")
		parts = strings.Split(path, "/")
	)

	if r.Method == "GET" {
	}

	if r.Method == "POST" {
		if len(parts) != 2 {
			fmt.Println("invalid action")
			http.Error(w, "Invalid action", http.StatusBadRequest)
			return
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Println("Failed to read request body:", err)
			http.Error(w, "Failed to read request body", http.StatusInternalServerError)
			return
		}
		p.producech <- &Message{
			Data:  body,
			Topic: parts[1],
		}
		var payloadMap map[string]interface{}
		if err := json.Unmarshal(body, &payloadMap); err != nil {
			fmt.Println("Failed to unmarshal JSON payload:", err)
			http.Error(w, "Failed to unmarshal JSON payload", http.StatusInternalServerError)
			return
		}
		response := map[string]interface{}{
			"status":  "success",
			"message": "Message forwarded successfully",
			"payload": payloadMap,
		}

		responseJSON, err := json.Marshal(response)
		if err != nil {
			fmt.Println("Failed to marshal JSON response:", err)
			http.Error(w, "Failed to marshal JSON response", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(responseJSON)
	}
}

func (p *HTTPProducer) Start() error {
	slog.Info("HTTP transport started", "port", p.listenAddr)
	return http.ListenAndServe(p.listenAddr, p)
}
