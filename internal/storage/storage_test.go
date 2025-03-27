package storage

import (
	"context"
	"crypto/internal/models"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestCryptoDB_Save(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	c := CryptoDB{db: db}

	t.Run("success", func(t *testing.T) {
		usage := &models.Data{
			Timestamp: time.Now(),
			Result: models.Result{
				Usdt_usd: models.Usdt_usd{
					Ask: &models.Ask{
						Price:     1.0,
						Quantity:  2.0,
						Timestamp: 3,
					},
					Bid: &models.Bid{
						Price:     4.0,
						Quantity:  5.0,
						Timestamp: 6,
					},
				},
			},
		}

		mock.ExpectExec("INSERT INTO crypto").WithArgs(usage.Timestamp, usage.Result.Usdt_usd.Ask.Price, usage.Result.Usdt_usd.Ask.Quantity, time.Unix(usage.Result.Usdt_usd.Ask.Timestamp, 0), usage.Result.Usdt_usd.Bid.Price, usage.Result.Usdt_usd.Bid.Quantity, time.Unix(usage.Result.Usdt_usd.Bid.Timestamp, 0)).WillReturnResult(sqlmock.NewResult(1, 1))

		err := c.Save(context.Background(), usage)
		assert.NoError(t, err)
	})

	t.Run("error", func(t *testing.T) {
		usage := &models.Data{
			Timestamp: time.Now(),
			Result: models.Result{
				Usdt_usd: models.Usdt_usd{
					Ask: &models.Ask{},
					Bid: &models.Bid{},
				},
			},
		}

		err := c.Save(context.Background(), usage)
		assert.Error(t, err)
	})
}
