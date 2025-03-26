package storage

import (
	"context"
	"crypto/internal/models"
	"database/sql"
	"fmt"

	"time"

	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

type Storage interface {
	Save(ctx context.Context, data *models.Data) error
}

type CryptoDB struct {
	db *sql.DB
}

func New(dbConfig string, log *zap.Logger) (*CryptoDB, error) {
	db, err := sql.Open("postgres", dbConfig)

	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		log.Error("failed to ping database", zap.Error(err))
	}

	return &CryptoDB{db: db}, nil
}

func (c *CryptoDB) Save(ctx context.Context, data *models.Data) error {
	query := "INSERT INTO crypto (time, ask_price, ask_volume, ask_time, bid_price, bid_volume, bid_time) VALUES ($1, $2, $3, $4, $5, $6, $7)"

	_, err := c.db.ExecContext(
		ctx,
		query,
		data.Timestamp,
		data.Result.Usdt_usd.Ask.Price,
		data.Result.Usdt_usd.Ask.Quantity,
		time.Unix(data.Result.Usdt_usd.Ask.Timestamp, 0),
		data.Result.Usdt_usd.Bid.Price,
		data.Result.Usdt_usd.Bid.Quantity,
		time.Unix(data.Result.Usdt_usd.Bid.Timestamp, 0),
	)

	if err != nil {
		return fmt.Errorf("failed to save data: %w", err)
	}
	return nil
}
