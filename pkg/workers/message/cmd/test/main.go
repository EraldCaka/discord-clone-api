package main

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"net/http"
)

func main() {
	url := "http://localhost:3003/publish"
	topics := []string{"hi", "how are you", "bye"}

	for i := 0; i < 1000; i++ {
		topic := topics[rand.Intn(len(topics))]
		payload := []byte(fmt.Sprintf("message_payload_%d", i))
		resp, err := http.Post(url+"/"+topic, "application/octet-stream", bytes.NewReader(payload))
		if err != nil {
			log.Fatal(err)
		}
		if resp.StatusCode != http.StatusOK {
			log.Fatal("status code is not 200")
		}
	}
}
