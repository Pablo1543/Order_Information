package nats

import (
	"Order_Information/internal/cache"
	"Order_Information/internal/db"
	"Order_Information/internal/models"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	stan "github.com/nats-io/stan.go"
)

const (
	ChannelName = "orders"
)

func StartSubscriber(clusterID, clientID, natsURL string, database *db.DB, memoryCache *cache.Cache) (stan.Conn, stan.Subscription, error) {
	sc, err := stan.Connect(clusterID, clientID, stan.NatsURL(natsURL))
	if err != nil {
		return nil, nil, err
	}

	sub, err := sc.Subscribe("orders", func(m *stan.Msg) {
		if err := handleMessage(database, memoryCache, m.Data); err != nil {
			log.Printf("Error handling message: %v", err)
			return
		}
		log.Printf("Order processed successfully")
	}, stan.DeliverAllAvailable())

	if err != nil {
		sc.Close()
		return nil, nil, err
	}
	return sc, sub, nil
}

func handleMessage(database *db.DB, memoryCache *cache.Cache, data []byte) error {
	var order models.Order
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&order); err != nil {
		return fmt.Errorf("invalid payload: %w", err)
	}
	if order.OrderUID == "" {
		return errors.New("empty order_uid")
	}
	// записываем в БД
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := database.SaveOrder(ctx, &order, data); err != nil {
		return fmt.Errorf("db save error: %w", err)
	}
	// обновляем кэш
	memoryCache.Set(order.OrderUID, json.RawMessage(data))
	log.Printf("order %s saved", order.OrderUID)
	return nil
}
