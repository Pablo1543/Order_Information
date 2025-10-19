package main

import (
	"Order_Information/internal/cache"
	"Order_Information/internal/db"
	"Order_Information/internal/nats"
	"context"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

type IndexData struct {
	Orders     []OrderPreview
	OrderCount int
}

type OrderPreview struct {
	OrderUID    string    `json:"order_uid"`
	TrackNumber string    `json:"track_number"`
	Entry       string    `json:"entry"`
	DateCreated time.Time `json:"date_created"`
}

// Функция для рендеринга главной страницы
func RenderIndexHTML(w http.ResponseWriter, data IndexData) {
	t, err := template.ParseFiles("pages/index.html")
	if err != nil {
		http.Error(w, "page error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if err := t.Execute(w, data); err != nil {
		http.Error(w, "page error: "+err.Error(), http.StatusInternalServerError)
	}
}

// Функция для рендеринга страницы заказа
func RenderOrderHTML(w http.ResponseWriter, data map[string]interface{}) {
	t, err := template.ParseFiles("pages/order.html")
	if err != nil {
		http.Error(w, "page error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if err := t.Execute(w, data); err != nil {
		http.Error(w, "page error: "+err.Error(), http.StatusInternalServerError)
	}
}

// Функция для получения превью заказов из кэша
func getOrderPreviews(cache *cache.Cache) ([]OrderPreview, error) {
	keys := cache.GetAllKeys()
	var orders []OrderPreview

	for _, key := range keys {
		raw, ok := cache.Get(key)
		if !ok {
			continue
		}

		var preview OrderPreview
		if err := json.Unmarshal(raw, &preview); err != nil {
			continue
		}
		orders = append(orders, preview)
	}

	return orders, nil
}

// Обработчик главной страницы
func indexHandler(w http.ResponseWriter, r *http.Request, memoryCache *cache.Cache) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	orders, err := getOrderPreviews(memoryCache)
	if err != nil {
		http.Error(w, "Error loading orders", http.StatusInternalServerError)
		return
	}

	data := IndexData{
		Orders:     orders,
		OrderCount: len(orders),
	}

	RenderIndexHTML(w, data)
}

func main() {
	pgUser := os.Getenv("DB_USER")
	pgPass := os.Getenv("DB_PASSWORD")
	pgHost := os.Getenv("DB_HOST")
	pgPort := os.Getenv("DB_PORT")
	pgName := os.Getenv("DB_NAME")

	pgDsn := "postgres://" + pgUser + ":" + pgPass + "@" + pgHost + ":" + pgPort + "/" + pgName

	natsURL := os.Getenv("NATS_URL")
	clusterID := os.Getenv("CLUSTER_ID")
	clientID := "order-service-1"

	ctx := context.Background()
	database, err := db.NewDB(ctx, pgDsn)
	if err != nil {
		log.Fatalf("db connect: %v", err)
	}
	defer database.Close()

	memoryCache := cache.NewCache()

	if err := memoryCache.LoadFromDB(database); err != nil {
		log.Printf("warning: couldn't load cache from db: %v", err)
	} else {
		log.Printf("cache restored from db: %d orders", memoryCache.Size())
	}

	sc, sub, err := nats.StartSubscriber(clusterID, clientID, natsURL, database, memoryCache)
	if err != nil {
		log.Fatalf("nats subscribe: %v", err)
	}
	defer func() {
		if sub != nil {
			sub.Close()
		}
		if sc != nil {
			sc.Close()
		}
	}()

	r := mux.NewRouter()

	// Главная страница
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		indexHandler(w, r, memoryCache)
	}).Methods("GET")

	// Страница заказа
	r.HandleFunc("/order/{id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]
		raw, ok := memoryCache.Get(id)
		if !ok {
			http.NotFound(w, r)
			return
		}
		if r.Header.Get("Accept") == "application/json" || r.URL.Query().Get("format") == "json" {
			w.Header().Set("Content-Type", "application/json")
			w.Write(raw)
			return
		}
		var pretty map[string]interface{}
		_ = json.Unmarshal(raw, &pretty)
		RenderOrderHTML(w, pretty)
	}).Methods("GET")

	srv := &http.Server{Addr: ":8080", Handler: r}

	go func() {
		log.Printf("http server listening on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("http listen: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
	log.Println("shutting down...")

	ctxShut, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctxShut)
}
