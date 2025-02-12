FROM golang:1.22.3-alpine
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main ./cmd
EXPOSE 50051
CMD ["./main"]
