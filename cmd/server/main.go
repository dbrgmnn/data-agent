package main

import (
	"context"
	"log"
	"monitoring/internal/config"
	initDB "monitoring/internal/db"
	"monitoring/internal/grpcserver"
	q "monitoring/internal/queue"
	"monitoring/proto"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// initialize database
	db, err := initDB.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// start RabbitMQ consumer in a goroutine
	rabbitURL := config.LoadConfig().RabbitURL

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		if err := q.StartMetricsConsumer(ctx, db, rabbitURL); err != nil {
			log.Fatal("Failed to start RabbitMQ consumer: ", err)
		}
	}()

	// start gRPC server in a goroutine
	go func() {
		lis, err := net.Listen("tcp", ":50051")
		if err != nil {
			log.Fatalf("Failed to listen: %v", err)
		}

		grpcServer := grpc.NewServer()
		proto.RegisterHostServiceServer(grpcServer, &grpcserver.HostService{DB: db})
		proto.RegisterMetricServiceServer(grpcServer, &grpcserver.MetricService{DB: db})

		// register reflection service on gRPC server
		reflection.Register(grpcServer)

		log.Println("gRPC server listening on :50051")

		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve gRPC server: %v", err)
		}
	}()

	// block main
	select {}
}
