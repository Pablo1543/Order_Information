package test

import (
	"Order_Information/internal/cache"
	"encoding/json"
	"fmt"
	"sync"
	"testing"
	"time"
)

// StressTestOrderCreation тестирует создание множества заказов
func TestStressOrderCreation(t *testing.T) {
	cache := cache.NewCache()

	concurrentUsers := 100
	ordersPerUser := 10
	var wg sync.WaitGroup
	wg.Add(concurrentUsers)

	start := time.Now()

	for i := 0; i < concurrentUsers; i++ {
		go func(userID int) {
			defer wg.Done()

			for j := 0; j < ordersPerUser; j++ {
				orderID := fmt.Sprintf("stress_%d_%d", userID, j)
				orderData := json.RawMessage(fmt.Sprintf(`{"order_uid": "%s", "user": %d}`, orderID, userID))
				cache.Set(orderID, orderData)
				time.Sleep(1 * time.Millisecond)
			}
		}(i)
	}

	wg.Wait()
	duration := time.Since(start)

	t.Logf("Created %d orders in %v", concurrentUsers*ordersPerUser, duration)
	t.Logf("Throughput: %.2f orders/second",
		float64(concurrentUsers*ordersPerUser)/duration.Seconds())
}
