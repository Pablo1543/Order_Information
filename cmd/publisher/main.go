package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	stan "github.com/nats-io/stan.go"
)

func connectWithRetry(clusterID, clientID, natsURL string, maxRetries int, delay time.Duration) (stan.Conn, error) {
	var sc stan.Conn
	var err error

	for i := 0; i < maxRetries; i++ {
		log.Printf("Attempting to connect to NATS Streaming (attempt %d/%d)...", i+1, maxRetries)

		sc, err = stan.Connect(clusterID, clientID, stan.NatsURL(natsURL))
		if err == nil {
			log.Println("Successfully connected to NATS Streaming!")
			return sc, nil
		}

		log.Printf("Connection failed: %v", err)
		if i < maxRetries-1 {
			log.Printf("Retrying in %v...", delay)
			time.Sleep(delay)
		}
	}

	return nil, fmt.Errorf("failed to connect to NATS Streaming after %d attempts: %w", maxRetries, err)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: publisher <file.json>")
		return
	}

	//Получаем параметры из переменных окружения
	clusterID := os.Getenv("CLUSTER_ID")
	if clusterID == "" {
		clusterID = "cluster1"
	}

	clientID := os.Getenv("CLIENT_ID")
	if clientID == "" {
		clientID = "publisher-1"
	}

	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		natsURL = "nats://localhost:4222"
	}

	//Подключаемся к NATS Streaming с повторными попытками
	sc, err := connectWithRetry(clusterID, clientID, natsURL, 10, 3*time.Second)
	if err != nil {
		log.Fatalf("NATS Streaming connection error: %v", err)
	}
	defer sc.Close()

	//Читаем файл
	data, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}

	//Публикуем сообщение
	if err := sc.Publish("orders", data); err != nil {
		log.Fatalf("Error publishing message: %v", err)
	}

	log.Println("Message published successfully!")
}
