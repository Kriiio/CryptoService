package service

import (
	"context"
	"crypto/internal/models"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type MockHTTPClient struct {
	Response *http.Response
	Err      error
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.Response, m.Err
}

type MockStorage struct {
	SaveFunc func(ctx context.Context, data *models.Data) error
}

func (m *MockStorage) Save(ctx context.Context, data *models.Data) error {
	return m.SaveFunc(ctx, data)
}

func TestConvertToAsks(t *testing.T) {
	tests := []struct {
		name    string
		rawData [][]interface{}
		wantErr bool
		wantAsk *models.Ask
	}{
		{
			name:    "empty rawData returns error",
			rawData: [][]interface{}{},
			wantErr: true,
		},
		{
			name: "rawData with invalid price format returns error",
			rawData: [][]interface{}{
				{"abc", "1", 123.45},
			},
			wantErr: true,
		},
		{
			name: "rawData with invalid quantity format returns error",
			rawData: [][]interface{}{
				{"123.45", "abc", 123.45},
			},
			wantErr: true,
		},
		{
			name: "rawData with valid format returns correct Ask object",
			rawData: [][]interface{}{
				{"123.45", "1", 123.45},
			},
			wantAsk: &models.Ask{
				Price:     123.45,
				Quantity:  1,
				Timestamp: 123,
			},
		},
		{
			name: "multiple items in rawData returns the last item's Ask object",
			rawData: [][]interface{}{
				{"123.45", "1", 123.45},
				{"234.56", "2", 234.56},
			},
			wantAsk: &models.Ask{
				Price:     234.56,
				Quantity:  2,
				Timestamp: 234,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ask, err := convertToAsks(tt.rawData)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantAsk, ask)
			}
		})
	}
}

func TestConvertToBids(t *testing.T) {
	tests := []struct {
		name    string
		rawData [][]interface{}
		wantErr bool
		wantBid *models.Bid
	}{
		{
			name:    "empty rawData returns nil and no error",
			rawData: [][]interface{}{},
			wantErr: false,
			wantBid: nil,
		},
		{
			name: "rawData with invalid price format returns error",
			rawData: [][]interface{}{
				{"abc", "1", 123.45},
			},
			wantErr: true,
			wantBid: nil,
		},
		{
			name: "rawData with invalid quantity format returns error",
			rawData: [][]interface{}{
				{"123.45", "abc", 123.45},
			},
			wantErr: true,
			wantBid: nil,
		},
		{
			name: "rawData with valid format returns correct Bid object",
			rawData: [][]interface{}{
				{"123.45", "1", 123.45},
			},
			wantErr: false,
			wantBid: &models.Bid{
				Price:     123.45,
				Quantity:  1,
				Timestamp: 123,
			},
		},
		{
			name: "rawData with multiple items returns the last item's Bid object",
			rawData: [][]interface{}{
				{"123.45", "1", 123.45},
				{"234.56", "2", 234.56},
			},
			wantErr: false,
			wantBid: &models.Bid{
				Price:     234.56,
				Quantity:  2,
				Timestamp: 234,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bid, err := convertToBids(tt.rawData)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantBid, bid)
			}
		})
	}
}

func TestFindRate(t *testing.T) {
	tests := []struct {
		name         string
		statusCode   int
		respBody     string
		expectedErr  error
		expectedData *models.Data
	}{
		{
			name:        "successful API call",
			statusCode:  http.StatusOK,
			respBody:    `{"result": {"USDTZUSD": {"asks": [["1.0", "1.0", 1643723911]], "bids": [["0.9", "1.0", 1643723900]]}}}`,
			expectedErr: nil,
			expectedData: &models.Data{
				Timestamp: time.Unix(1643723911, 0),
				Result: models.Result{
					Usdt_usd: models.Usdt_usd{
						Ask: &models.Ask{Price: 1.0, Quantity: 1.0, Timestamp: 1643723911},
						Bid: &models.Bid{Price: 0.9, Quantity: 1.0, Timestamp: 1643723900},
					},
				},
			},
		},
		{
			name:         "failed API call",
			statusCode:   http.StatusInternalServerError,
			respBody:     "",
			expectedErr:  errors.New("service.findRate: failed to decode EOF"),
			expectedData: nil,
		},
		{
			name:         "failed request creation",
			statusCode:   0,
			respBody:     "",
			expectedErr:  errors.New("service.findRate: failed to decode EOF"),
			expectedData: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Создаем mock-клиент
			mockClient := &MockHTTPClient{
				Response: &http.Response{
					StatusCode: test.statusCode,
					Body:       io.NopCloser(strings.NewReader(test.respBody)),
				},
				Err: nil,
			}

			// Вызываем функцию с mock-клиентом
			data, err := findRate(mockClient)

			// Проверяем ошибки
			if test.expectedErr != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), test.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}

			// Проверяем данные
			if test.expectedData != nil {
				assert.NotNil(t, data)
				assert.Equal(t, test.expectedData.Result.Usdt_usd.Ask.Price, data.Result.Usdt_usd.Ask.Price)
				assert.Equal(t, test.expectedData.Result.Usdt_usd.Ask.Quantity, data.Result.Usdt_usd.Ask.Quantity)
				assert.Equal(t, test.expectedData.Result.Usdt_usd.Ask.Timestamp, data.Result.Usdt_usd.Ask.Timestamp)
				assert.Equal(t, test.expectedData.Result.Usdt_usd.Bid.Price, data.Result.Usdt_usd.Bid.Price)
				assert.Equal(t, test.expectedData.Result.Usdt_usd.Bid.Quantity, data.Result.Usdt_usd.Bid.Quantity)
				assert.Equal(t, test.expectedData.Result.Usdt_usd.Bid.Timestamp, data.Result.Usdt_usd.Bid.Timestamp)
			} else {
				assert.Nil(t, data)
			}
		})
	}
}

func TestGetData(t *testing.T) {
	tests := []struct {
		name         string
		statusCode   int
		respBody     string
		findRateErr  error
		saveErr      error
		expectedErr  error
		expectedData *models.Data
	}{
		{
			name:        "successful API call and data save",
			findRateErr: nil,
			statusCode:  http.StatusOK,
			respBody:    `{"result": {"USDTZUSD": {"asks": [["1.0", "1.0", 1643723911]], "bids": [["0.9", "1.0", 1643723900]]}}}`,
			saveErr:     nil,
			expectedErr: nil,
			expectedData: &models.Data{
				Timestamp: time.Unix(1643723911, 0),
				Result: models.Result{
					Usdt_usd: models.Usdt_usd{
						Ask: &models.Ask{Price: 1.0, Quantity: 1.0, Timestamp: 1643723911},
						Bid: &models.Bid{Price: 0.9, Quantity: 1.0, Timestamp: 1643723900},
					},
				},
			},
		},
		{
			name:         "failed data save",
			statusCode:   http.StatusOK,
			respBody:     `{"result": {"USDTZUSD": {"asks": [["1.0", "1.0", 1643723911]], "bids": [["0.9", "1.0", 1643723900]]}}}`,
			findRateErr:  nil,
			saveErr:      errors.New("save error"),
			expectedErr:  errors.New("save error"),
			expectedData: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Create a mock logger
			logger := zap.NewNop()

			// Create a mock database
			db := &MockStorage{
				SaveFunc: func(ctx context.Context, data *models.Data) error {
					return test.saveErr
				},
			}

			// Create a mock HTTP client
			mockClient := &MockHTTPClient{
				Response: &http.Response{
					StatusCode: test.statusCode,
					Body:       io.NopCloser(strings.NewReader(test.respBody)),
				},
				Err: nil,
			}

			// Create a service instance
			s := &ServiceImpl{
				log:    logger,
				db:     db,
				client: mockClient,
			}

			// Call the GetData function
			data, err := s.GetData(context.Background())

			// Assert the expected error and data

			if test.expectedData != nil {
				assert.NotNil(t, data)
				assert.Equal(t, err, test.expectedErr)
				assert.Equal(t, test.expectedData.Result.Usdt_usd.Ask.Price, data.Result.Usdt_usd.Ask.Price)
				assert.Equal(t, test.expectedData.Result.Usdt_usd.Ask.Quantity, data.Result.Usdt_usd.Ask.Quantity)
				assert.Equal(t, test.expectedData.Result.Usdt_usd.Ask.Timestamp, data.Result.Usdt_usd.Ask.Timestamp)
				assert.Equal(t, test.expectedData.Result.Usdt_usd.Bid.Price, data.Result.Usdt_usd.Bid.Price)
				assert.Equal(t, test.expectedData.Result.Usdt_usd.Bid.Quantity, data.Result.Usdt_usd.Bid.Quantity)
				assert.Equal(t, test.expectedData.Result.Usdt_usd.Bid.Timestamp, data.Result.Usdt_usd.Bid.Timestamp)
			} else {
				assert.Nil(t, data)
			}
		})
	}
}
