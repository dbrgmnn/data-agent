package main

import (
	"data-agent/internal/config"
	dataBase "data-agent/internal/db"
	"data-agent/internal/grpcserver"
	"data-agent/proto"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// main function to start the gRPC server
func main() {
	// add prefix for logs
	log.SetPrefix("[api] ")
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// initialize database
	db, err := dataBase.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// load configuration
	cfg := config.LoadConfig()
	grpcPort := cfg.GRPCPort

	// start gRPC server
	lis, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	proto.RegisterHostServiceServer(grpcServer, &grpcserver.HostService{DB: db})
	proto.RegisterMetricServiceServer(grpcServer, &grpcserver.MetricService{DB: db})
	reflection.Register(grpcServer)

	go func() {
		log.Printf("gRPC server started on %s", grpcPort)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve gRPC server: %v", err)
		}
	}()

	// handle termination signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// wait for termination signal
	<-stop
	log.Println("Stopping gRPC server...")
	grpcServer.GracefulStop()
}
