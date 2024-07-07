# 🚀 Golang gRPC User Service with Search

This project implements a Golang gRPC service for managing user details with search functionality using an SQLite database.

## 🎯 Objective

The objective of this project is to develop a gRPC service that provides specific functionalities for managing user details and includes a search capability. The primary objectives are as follows:

- Simulate a database with user details.
- Create gRPC endpoints for fetching user details by ID, fetching users by a list of IDs, and searching users based on specific criteria.
- Ensure code quality, design patterns, and comprehensive unit tests for reliability.

## 🛠️ Technologies Used

- Backend: Golang
- Database: SQLite
- gRPC Framework
- Docker (optional for deployment)

## 🏗️ Project Structure

The project structure is organized as follows:

```
golang-grpc-user-services/
├── cmd/
│   └── main.go         # Main entry point of the gRPC server
├── internal/
│   ├── database/
│   │   └── db.go       # Database initialization and operations
│   ├── server/
│   │   ├── server.go   # gRPC server implementation
│   │   └── logger.go   # Centralized logger setup
│   └── logger/
│       └── logger.go   # Logger package for centralized logging
├── pb/
│   ├── user.proto      # Protobuf file defining user service APIs
│   └── user.pb.go      # Generated Go code from user.proto
├── db/
│   └── sahil89.db      # SQLite database file
└── go.mod              # Go module file
```


## 🚀 Getting Started

### 📋 Prerequisites

- Go installed on your machine.
- Docker (optional for deployment).
- Clone this repository:

  ```bash
  git clone https://github.com/Sahil-4555/Golang-gRPC-User-Service-with-Search.git
  ````

## 🏃‍♂️ Running the Application
Running Locally
Start the gRPC server:
```bash
go run cmd/main.go
```

## 📝 Implementation Details

- **GetUser**: Fetches user details by ID.
- **GetUsers**: Retrieves user details by a list of IDs.
- **SearchUsers**: Searches users based on criteria such as city, phone number, and marital status.

## 🧪 Testing

Unit tests cover each gRPC method and various edge cases. To run tests, use:

```bash
go test ./...
```

### Using grpcurl

grpcurl is a command-line tool that allows you to interact with gRPC servers. Here are a few examples of how to use grpcurl to test your gRPC methods:

1. **GetUser**:
   Fetch user details by ID.
```bash
grpcurl -plaintext -d '{"id": 1}' localhost:50051 user.UserService/GetUser   
```
2. **GetUsers**:
    Fetch users details by ID.
```bash
grpcurl -plaintext -d '{"ids": [1,3,5,9]}' localhost:50051 user.UserService/GetUsers   
```
3. **SearchUser**:
    Search User by Fields.
```bash
grpcurl -plaintext -d '{"phone": "555", "married": true}' localhost:50051 user.UserService/SearchUsers
```
    



