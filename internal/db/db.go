package db

import (
	"Order_Information/internal/models"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	pool *pgxpool.Pool
}

func NewDB(ctx context.Context, dsn string) (*DB, error) {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}
	return &DB{pool: pool}, nil
}

func (db *DB) Close() {
	if db.pool != nil {
		db.pool.Close()
	}
}

func (db *DB) SaveOrder(ctx context.Context, order *models.Order, raw []byte) error {
	_, err := db.pool.Exec(ctx, `
        INSERT INTO orders (order_uid, track_number, customer_id, date_created, raw)
        VALUES ($1, $2, $3, $4, $5)
        ON CONFLICT (order_uid) DO UPDATE
        SET track_number = EXCLUDED.track_number,
            customer_id = EXCLUDED.customer_id,
            date_created = EXCLUDED.date_created,
            raw = EXCLUDED.raw
    `, order.OrderUID, order.TrackNumber, order.CustomerID, order.DateCreated, raw)
	return err
}

func (db *DB) GetOrderRaw(ctx context.Context, orderUID string) ([]byte, error) {
	var raw []byte
	row := db.pool.QueryRow(ctx, `SELECT raw FROM orders WHERE order_uid=$1`, orderUID)
	if err := row.Scan(&raw); err != nil {
		return nil, err
	}
	return raw, nil
}

func (db *DB) LoadAllOrders(ctx context.Context) (map[string][]byte, error) {
	rows, err := db.pool.Query(ctx, `SELECT order_uid, raw FROM orders`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := make(map[string][]byte)
	for rows.Next() {
		var uid string
		var raw []byte
		if err := rows.Scan(&uid, &raw); err != nil {
			return nil, err
		}
		res[uid] = raw
	}
	return res, nil
}
