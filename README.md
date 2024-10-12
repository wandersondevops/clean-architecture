# Clean Architecture Orders System

This project is a **Clean Architecture** implementation for managing orders. It supports three types of APIs for order management:

- **REST** API on port `8000`
- **gRPC** service on port `50051`
- **GraphQL** API on port `8080`

The project allows creating and listing orders through these three interfaces. The backend is developed in Go, and it uses MySQL as the database.

## Features

- **Create and List Orders**: You can create and list orders using REST, gRPC, and GraphQL.
- **Clean Architecture**: The project follows the principles of Clean Architecture, ensuring decoupling between layers.
- **Docker**: The project is fully containerized using Docker Compose.

## Technologies Used

- **Go**: The main programming language for the project.
- **gRPC**: For handling high-performance, strongly typed communications.
- **GraphQL**: For flexible and powerful API querying.
- **REST**: For simple, widely-used API integration.
- **MySQL**: The database used to store order data.
- **Docker**: For containerization and environment setup.

## Prerequisites

Before you start, make sure you have the following installed:

- **Go** (v1.18 or later): [Download Go](https://golang.org/dl/)
- **Docker**: [Install Docker](https://docs.docker.com/get-docker/)
- **Docker Compose**: [Install Docker Compose](https://docs.docker.com/compose/install/)

## Setup Instructions

### 1. Clone the Repository

```bash
git clone https://github.com/wandersondevops/clean-architecture.git
cd clean-architecture-orders-system
```

### 2. Install Go Dependencies

```bash
go mod tidy
```

### 4. Install gRPC and Protobuf Go Plugins
```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2.0
```

### 5. Make sure the protoc-gen-go and protoc-gen-go-grpc plugins are in your PATH:

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

### 6. Generate gRPC and Protobuf Files
Run the following command to generate the necessary Go code from the proto files:

```bash
protoc --go_out=. --go-grpc_out=. internal/infra/grpc/protofiles/order.proto
```

### Build and Run the Project with Docker Compose
To build and start the containers:

```bash
docker compose up --build
```

#### This will:

- Start a MySQL database on port 3306
- Start the REST API on port 8000
- Start the gRPC server on port 50051
- Start the GraphQL API on port 8080

### Applying Database Migrations
To run the MySQL database migrations, use the migrate tool:

```bash
migrate -path ./internal/infra/database/migrations -database "mysql://root:root@tcp(localhost:3306)/orders" up
```

### Testing the Web Servers

#### REST API (Port 8000)
You can interact with the REST API using curl or Postman.

#### Create an Order (POST /order):

```bash
curl -X POST http://localhost:8000/order -H "Content-Type: application/json" -d '{
    "id": "order1",
    "price": 100.5,
    "tax": 10.5
}'
```

#### List Orders (GET /orders):

```bash
curl http://localhost:8000/orders
gRPC Service (Port 50051)
```

#### You can test the gRPC service using grpcurl:

```bash
grpcurl -plaintext -d '{}' localhost:50051 pb.OrderService/ListOrders
```

#### GraphQL API (Port 8080)

You can access the GraphQL playground at http://localhost:8080.

- List Orders:

query {
  listOrders {
    id
    price
    tax
    finalPrice
  }
}

- Create an Order:

mutation {
  createOrder(input: { id: "order2", price: 150.75, tax: 15.50 }) {
    id
    price
    tax
    finalPrice
  }
}