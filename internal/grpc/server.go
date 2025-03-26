package server

import (
	"context"
	"crypto/internal/service"

	cryptoservicev1 "github.com/Kriiio/proto/gen/go/usdt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type serverAPI struct {
	cryptoservicev1.UnimplementedCryptoproviderServer
	service service.Service
}

func Register(gRPC *grpc.Server, service service.Service) {
	cryptoservicev1.RegisterCryptoproviderServer(gRPC, &serverAPI{service: service})
}

func (s *serverAPI) GetRates(ctx context.Context, req *cryptoservicev1.Request) (*cryptoservicev1.Response, error) {
	data, err := s.service.GetData(ctx)

	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get data")
	}

	ask := data.Result.Usdt_usd.Ask
	bid := data.Result.Usdt_usd.Bid

	var response cryptoservicev1.Response
	response.Asks = append(response.Asks, &cryptoservicev1.Ask{
		Price:     ask.Price,
		Volume:    ask.Quantity,
		Timestamp: ask.Timestamp,
	})

	response.Bids = append(response.Bids, &cryptoservicev1.Bid{
		Price:     bid.Price,
		Volume:    bid.Quantity,
		Timestamp: bid.Timestamp,
	})

	response.Timestamp = data.Timestamp.Unix()

	return &response, nil
}
