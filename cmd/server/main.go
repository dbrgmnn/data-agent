package main

import (
	"log"
	"monitoring/internal/config"
	"monitoring/internal/grpcserver"
	pb "monitoring/internal/grpcserver/gen"
	"monitoring/internal/server"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// initialize database
	db, err := server.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// start RabbitMQ consumer in a goroutine
	rabbitURL := config.LoadConfig().RabbitURL
	go func() {
		if err := server.StartMetricsConsumer(db, rabbitURL); err != nil {
			log.Fatal(err)
		}
	}()

	// start gRPC server in a goroutine
	go func() {
		lis, err := net.Listen("tcp", ":50051")
		if err != nil {
			log.Fatalf("Failed to listen: %v", err)
		}

		grpcServer := grpc.NewServer()
		metricsServer := grpcserver.NewMetricsServer(db)
		pb.RegisterMetricsServiceServer(grpcServer, metricsServer)

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
