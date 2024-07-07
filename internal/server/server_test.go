package server

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"
	"testing"

	"golang-grcp-user-services/internal/database"
	pb "golang-grcp-user-services/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

var (
	testDBFile = "../../db/test_db.db"
)

func setup() (*sql.DB, error) {
	db, err := database.InitDB(testDBFile)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}
	return db, nil
}

func bufListener() (*bufconn.Listener, error) {
	listener := bufconn.Listen(1024 * 1024)
	return listener, nil
}

func TestGetUser(t *testing.T) {
	db, err := setup()
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}
	defer db.Close()

	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, &server{db: db})

	lis, err := bufListener()
	if err != nil {
		t.Fatalf("Failed to create buf listener: %v", err)
	}
	defer lis.Close()

	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()
	defer s.Stop()

	conn, err := grpc.DialContext(context.Background(), "bufnet", grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	client := pb.NewUserServiceClient(conn)

	// Test case 1: Valid user ID
	req := &pb.UserID{Id: 1}
	user, err := client.GetUser(context.Background(), req)
	if err != nil {
		t.Fatalf("GetUser returned an error: %v", err)
	}
	if user == nil {
		t.Fatalf("Expected a user, got nil")
	}
	if user.Id != 1 {
		t.Errorf("Expected user ID to be 1, got %d", user.Id)
	}

	// Test case 2: Non-existing user ID
	req = &pb.UserID{Id: 1000}
	user, err = client.GetUser(context.Background(), req)
	if user != nil || err == nil {
		t.Fatalf("Expected user to be nil and error to be non-nil")
	}
}

// TestGetUsers tests the GetUsers RPC method.
func TestGetUsers(t *testing.T) {
	db, err := setup()
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}
	defer db.Close()

	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, &server{db: db})

	lis, err := bufListener()
	if err != nil {
		t.Fatalf("Failed to create buf listener: %v", err)
	}
	defer lis.Close()

	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()
	defer s.Stop()

	conn, err := grpc.DialContext(context.Background(), "bufnet", grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	client := pb.NewUserServiceClient(conn)

	// Test case: Valid user IDs
	req := &pb.UserIDs{Ids: []int32{1, 2, 3}}
	users, err := client.GetUsers(context.Background(), req)
	if err != nil {
		t.Fatalf("GetUsers returned an error: %v", err)
	}
	if len(users.Users) != 3 {
		t.Fatalf("Expected 3 users, got %d", len(users.Users))
	}
	for _, user := range users.Users {
		if user.Id != 1 && user.Id != 2 && user.Id != 3 {
			t.Errorf("Unexpected user ID: %d", user.Id)
		}
	}
}

// TestSearchUsers tests the SearchUsers RPC method.
func TestSearchUsers(t *testing.T) {
	db, err := setup()
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}
	defer db.Close()

	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, &server{db: db})

	lis, err := bufListener()
	if err != nil {
		t.Fatalf("Failed to create buf listener: %v", err)
	}
	defer lis.Close()

	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()
	defer s.Stop()

	conn, err := grpc.DialContext(context.Background(), "bufnet", grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	client := pb.NewUserServiceClient(conn)

	// Test case: Search by city "Mumbai"
	req := &pb.SearchCriteria{City: "Mumbai"}
	users, err := client.SearchUsers(context.Background(), req)
	if err != nil {
		t.Fatalf("SearchUsers returned an error: %v", err)
	}
	if len(users.Users) == 0 {
		t.Fatalf("Expected users in Mumbai, got none")
	}
	for _, user := range users.Users {
		if user.City != "Mumbai" {
			t.Errorf("Expected user in Mumbai, got user in %s", user.City)
		}
	}

	// Test case: Search by phone number 1234567890
	req = &pb.SearchCriteria{Phone: 1234567890}
	users, err = client.SearchUsers(context.Background(), req)
	if err != nil {
		t.Fatalf("SearchUsers returned an error: %v", err)
	}
	if len(users.Users) == 0 {
		t.Fatalf("Expected users with phone number 1234567890, got none")
	}
	for _, user := range users.Users {
		if user.Phone != 1234567890 {
			t.Errorf("Expected user with phone number 1234567890, got user with phone number %d", user.Phone)
		}
	}

	// Test case: Search by married status
	req = &pb.SearchCriteria{MarriedCriteria: &pb.SearchCriteria_Married{Married: true}}
	users, err = client.SearchUsers(context.Background(), req)
	if err != nil {
		t.Fatalf("SearchUsers returned an error: %v", err)
	}
	if len(users.Users) == 0 {
		t.Fatalf("Expected married users, got none")
	}
	for _, user := range users.Users {
		if !user.Married {
			t.Errorf("Expected married user, got unmarried user")
		}
	}

	// Test case: Search by city "Mumbai" and married status
	req = &pb.SearchCriteria{City: "Mumbai", MarriedCriteria: &pb.SearchCriteria_Married{Married: true}}
	users, err = client.SearchUsers(context.Background(), req)
	if err != nil {
		t.Fatalf("SearchUsers returned an error: %v", err)
	}
	if len(users.Users) == 0 {
		t.Fatalf("Expected married users in Mumbai, got none")
	}
	for _, user := range users.Users {
		if user.City != "Mumbai" {
			t.Errorf("Expected user in Mumbai, got user in %s", user.City)
		}
		if !user.Married {
			t.Errorf("Expected married user, got unmarried user")
		}
	}

	// Test case: Search by phone number 1234567890 and married status
	req = &pb.SearchCriteria{Phone: 1234567890, MarriedCriteria: &pb.SearchCriteria_Married{Married: true}}
	users, err = client.SearchUsers(context.Background(), req)
	if err != nil {
		t.Fatalf("SearchUsers returned an error: %v", err)
	}
	if len(users.Users) == 0 {
		t.Fatalf("Expected users with phone number 1234567890 and married, got none")
	}
	for _, user := range users.Users {
		if user.Phone != 1234567890 {
			t.Errorf("Expected user with phone number 1234567890, got user with phone number %d", user.Phone)
		}
		if !user.Married {
			t.Errorf("Expected married user, got unmarried user")
		}
	}

}

func cleanup() {
	err := os.Remove(testDBFile)
	if err != nil {
		log.Printf("Failed to delete test database file: %v", err)
	}
}

func TestMain(m *testing.M) {
	defer cleanup()
	code := m.Run()
	os.Exit(code)
}
