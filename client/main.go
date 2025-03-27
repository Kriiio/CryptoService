package main

import (
	"context"
	"fmt"
	"log"

	cryptoservicev1 "github.com/Kriiio/proto/gen/go/usdt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

func main() {
	conn, err := grpc.DialContext(context.Background(), "crypto-grpc-service:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Ошибка при подключении к серверу: %v", err)
	}
	fmt.Println("Подключено")
	defer conn.Close()

	client := cryptoservicev1.NewCryptoproviderClient(conn)

	res, err := client.GetRates(context.Background(), &emptypb.Empty{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res)

	ping, err := client.Ping(context.Background(), &emptypb.Empty{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(ping.Pong)
}
