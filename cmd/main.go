package main

import (
	"net"

	"golang-grcp-user-services/internal/database"
	"golang-grcp-user-services/internal/logger"
	"golang-grcp-user-services/internal/server"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		logger.ErrorLogger.Fatalf("Failed to listen: %v", err)
	}

	db, err := database.InitDB("db/user.db")
	if err != nil {
		logger.ErrorLogger.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	s := server.NewServer(db)

	logger.InfoLogger.Printf("Server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		logger.ErrorLogger.Fatalf("Failed to serve: %v", err)
	}
}
