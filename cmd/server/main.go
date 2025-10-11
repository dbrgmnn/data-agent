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
	"os"
	"os/signal"
	"syscall"

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

	// load configuration
	cfg := config.LoadConfig()
	rabbitURL := cfg.RabbitURL
	grpcPort := cfg.GRPCPort

	// create a context that is canceled on exit
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// start RabbitMQ consumer
	go func() {
		if err := q.StartMetricsConsumer(ctx, db, rabbitURL); err != nil {
			log.Printf("Failed to start RabbitMQ consumer: %v", err)
		}
	}()

	// start gRPC server
	lis, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	proto.RegisterHostServiceServer(grpcServer, &grpcserver.HostService{DB: db})
	proto.RegisterMetricServiceServer(grpcServer, &grpcserver.MetricService{DB: db})

	// register reflection service on gRPC server
	reflection.Register(grpcServer)

	go func() {
		log.Printf("gRPC server listening on %s", grpcPort)

		if err := grpcServer.Serve(lis); err != nil {
			log.Printf("Failed to serve gRPC server: %v", err)
			cancel()
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
	log.Println("Shutting down server...")
	cancel()
	grpcServer.GracefulStop()
	log.Println("Server stopped.")
}
