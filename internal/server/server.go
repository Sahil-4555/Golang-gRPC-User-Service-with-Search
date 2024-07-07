package server

import (
	"context"
	"database/sql"
	"fmt"

	"golang-grcp-user-services/internal/logger"
	pb "golang-grcp-user-services/pb"

	_ "github.com/mattn/go-sqlite3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	pb.UnimplementedUserServiceServer
	db *sql.DB
}

func (s *server) GetUser(ctx context.Context, in *pb.UserID) (*pb.User, error) {
	var user pb.User
	err := s.db.QueryRow("SELECT id, fname, city, phone, height, married FROM users WHERE id = ?", in.Id).Scan(&user.Id, &user.Fname, &user.City, &user.Phone, &user.Height, &user.Married)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.InfoLogger.Printf("User with ID %d not found", in.Id)
			return nil, fmt.Errorf("user not found")
		}
		logger.ErrorLogger.Printf("Failed to fetch user with ID %d: %v", in.Id, err)
		return nil, fmt.Errorf("failed to fetch user: %v", err)
	}
	logger.InfoLogger.Printf("Successfully fetched user with ID %d", in.Id)
	return &user, nil
}

func (s *server) GetUsers(ctx context.Context, in *pb.UserIDs) (*pb.Users, error) {
	var users []*pb.User
	for _, id := range in.Ids {
		var user pb.User
		err := s.db.QueryRow("SELECT id, fname, city, phone, height, married FROM users WHERE id = ?", id).Scan(&user.Id, &user.Fname, &user.City, &user.Phone, &user.Height, &user.Married)
		if err != nil {
			if err == sql.ErrNoRows {
				logger.InfoLogger.Printf("User with ID %d not found", id)
				continue
			}
			logger.ErrorLogger.Printf("Failed to fetch user with ID %d: %v", id, err)
			return nil, fmt.Errorf("failed to fetch user with ID %d: %v", id, err)
		}
		users = append(users, &user)
		logger.InfoLogger.Printf("Successfully fetched user with ID %d", id)
	}
	return &pb.Users{Users: users}, nil
}

func (s *server) SearchUsers(ctx context.Context, in *pb.SearchCriteria) (*pb.Users, error) {
	var users []*pb.User
	query := "SELECT id, fname, city, phone, height, married FROM users WHERE 1=1"
	var args []interface{}

	if in.City != "" {
		query += " AND city LIKE ?"
		args = append(args, "%"+in.City+"%")
	}
	if in.Phone != 0 {
		phoneStr := fmt.Sprintf("%d", in.Phone)
		query += " AND CAST(phone AS TEXT) LIKE ?"
		args = append(args, "%"+phoneStr+"%")
	}
	if marriedCriteria, ok := in.MarriedCriteria.(*pb.SearchCriteria_Married); ok {
		query += " AND married = ?"
		args = append(args, marriedCriteria.Married)
	}

	rows, err := s.db.Query(query, args...)
	if err != nil {
		logger.ErrorLogger.Printf("Failed to execute query: %v", err)
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var user pb.User
		err := rows.Scan(&user.Id, &user.Fname, &user.City, &user.Phone, &user.Height, &user.Married)
		if err != nil {
			logger.ErrorLogger.Printf("Failed to scan row: %v", err)
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}
		users = append(users, &user)
	}
	if err := rows.Err(); err != nil {
		logger.ErrorLogger.Printf("Error iterating rows: %v", err)
		return nil, fmt.Errorf("error iterating rows: %v", err)
	}

	logger.InfoLogger.Println("Search query executed successfully")
	return &pb.Users{Users: users}, nil
}

func NewServer(db *sql.DB) *grpc.Server {
	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, &server{db: db})
	reflection.Register(s)
	return s
}
