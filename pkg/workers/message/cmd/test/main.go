package main

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

func main() {
	url := "http://localhost:3003/publish"
	topics := []string{"hi", "how are you", "bye"}
	totalRequests := 100000

	startTime := time.Now()

	for i := 0; i < totalRequests; i++ {
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

	elapsedTime := time.Since(startTime)
	requestsPerSecond := float64(totalRequests) / elapsedTime.Seconds()

	fmt.Printf("Total Requests: %d\n", totalRequests)
	fmt.Printf("Elapsed Time: %s\n", elapsedTime)
	fmt.Printf("Requests Per Second: %.2f\n", requestsPerSecond)
}
