package main

import (
	"context"
	"distributed-queue/config"
	"distributed-queue/internal/handler/rest"
	"distributed-queue/internal/pkg/filestorage"
	"distributed-queue/internal/repository/filerepository"
	"distributed-queue/internal/service"
	"distributed-queue/internal/transport/rpc"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	configPath := flag.String("config", "../../config/config.json", "Path to config file")
	flag.Parse()

	// Load configuration
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Setup repositories
	// queueRepo := memory.NewQueueRepository()

	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get current directory: %v", err)
	}

	queueFileStorage, err := filestorage.NewFileStorage(currentDir + "/data/" + cfg.NodeID + "/queue_storage")
	if err != nil {
		log.Fatalf("Failed to create queue file storage: %v", err)
	}

	clientFileStorage, err := filestorage.NewFileStorage(currentDir + "/data/" + cfg.NodeID + "/client_storage")
	if err != nil {
		log.Fatalf("Failed to create client file storage: %v", err)
	}

	queueRepo := filerepository.NewQueueRepository(queueFileStorage, clientFileStorage)

	// Setup RPC client/server
	rpcClient := rpc.NewClient(cfg.NodeTimeout)
	rpcServer := rpc.NewServer()

	// Setup services
	nodeService := service.NewNodeService(cfg, queueRepo, rpcClient)
	queueService := service.NewQueueService(cfg, queueRepo, rpcClient, nodeService)

	// Start node health check
	go nodeService.StartHealthCheck(cfg.HealthCheckInterval)

	// Start RPC server
	go func() {
		if err := rpcServer.Start(cfg.RPCPort, queueService); err != nil {
			log.Fatalf("Failed to start RPC server: %v", err)
		}
	}()

	// Setup REST API handler
	handler := rest.NewHandler(queueService)
	server := &http.Server{
		Addr:    ":" + cfg.HTTPPort,
		Handler: handler.SetupRoutes(),
	}

	// Start HTTP server
	go func() {
		log.Printf("Starting HTTP server on port %s", cfg.HTTPPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start HTTP server: %v", err)
		}
	}()

	// Wait for termination signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Create context with timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Gracefully shutdown HTTP server
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	// Gracefully shutdown RPC server
	rpcServer.Stop()

	log.Println("Server exited properly")
}
