package service

import (
	"context"
	"crypto/internal/models"
	"crypto/internal/storage"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"go.uber.org/zap"
)

type Service interface {
	GetData(ctx context.Context) (*models.Data, error)
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type ServiceImpl struct {
	log    *zap.Logger
	db     storage.Storage
	client HTTPClient
}

func New(log *zap.Logger, db storage.Storage, client *http.Client) *ServiceImpl {
	return &ServiceImpl{
		log:    log,
		db:     db,
		client: client,
	}
}

func (s *ServiceImpl) GetData(ctx context.Context) (*models.Data, error) {
	data, err := findRate(s.client)

	if err != nil {
		s.log.Error("failed to find rate", zap.Error(err))
		return nil, err
	}

	if err := s.db.Save(ctx, data); err != nil {
		s.log.Error("failed to save data", zap.Error(err))
		return nil, err
	}

	return data, nil
}

func findRate(client HTTPClient) (*models.Data, error) {
	const op = "service.findRate"

	data := &models.Data{}
	data.Timestamp = time.Now()

	url := "https://api.kraken.com/0/public/Depth?pair=USDTUSD&count=1"
	method := "GET"

	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		return nil, errors.New("service.findRate: failed to create request")
	}
	req.Header.Add("Accept", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return nil, errors.New("service.findRate: failed to send request")
	}

	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to decode %w", op, err)
	}

	data.Result.Usdt_usd.Ask, err = convertToAsks(data.Result.Usdt_usd.RawAsks)
	if err != nil {
		return nil, err
	}
	data.Result.Usdt_usd.Bid, err = convertToBids(data.Result.Usdt_usd.RawBids)
	if err != nil {
		return nil, err
	}

	return data, nil

}

func convertToAsks(rawData [][]interface{}) (*models.Ask, error) {
	const op = "service.convertToAsks"

	var ask *models.Ask

	if len(rawData) == 0 {
		return nil, fmt.Errorf("%s:raw data is empty", op)
	}

	for _, item := range rawData {
		price := item[0].(string)
		priceFloat, err := strconv.ParseFloat(price, 64)
		if err != nil {
			return nil, fmt.Errorf("%s: cannot parse price %w", op, err)
		}

		quantity := item[1].(string)
		quantityFloat, err := strconv.ParseFloat(quantity, 64)
		if err != nil {
			return nil, fmt.Errorf("%s: cannot parse quantity %w", op, err)
		}

		timestamp := int64(item[2].(float64))

		ask = &models.Ask{
			Price:     priceFloat,
			Quantity:  quantityFloat,
			Timestamp: timestamp,
		}
	}
	return ask, nil
}

func convertToBids(rawData [][]interface{}) (*models.Bid, error) {
	const op = "service.convertToBids"
	var bid *models.Bid
	for _, item := range rawData {
		price := item[0].(string)
		priceFloat, err := strconv.ParseFloat(price, 64)
		if err != nil {
			return nil, fmt.Errorf("%s: cannot parse price %w", op, err)
		}

		quantity := item[1].(string)
		quantityFloat, err := strconv.ParseFloat(quantity, 64)
		if err != nil {
			return nil, fmt.Errorf("%s: cannot parse quantity %w", op, err)
		}

		timestamp := int64(item[2].(float64))

		bid = &models.Bid{
			Price:     priceFloat,
			Quantity:  quantityFloat,
			Timestamp: timestamp,
		}
	}
	return bid, nil
}
