package main

import (
	"database/sql"
	"fmt"
	"net"
	"net/http"
	"sync"

	graphql_handler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/devfullcycle/20-CleanArch/configs"
	"github.com/devfullcycle/20-CleanArch/internal/event/handler"
	"github.com/devfullcycle/20-CleanArch/internal/infra/database"
	"github.com/devfullcycle/20-CleanArch/internal/infra/graph"
	"github.com/devfullcycle/20-CleanArch/internal/infra/grpc/pb"
	"github.com/devfullcycle/20-CleanArch/internal/infra/grpc/service"
	"github.com/devfullcycle/20-CleanArch/internal/infra/web"
	"github.com/devfullcycle/20-CleanArch/internal/usecase"
	"github.com/devfullcycle/20-CleanArch/pkg/events"
	"github.com/streadway/amqp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	// mysql
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	configs, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	db, err := sql.Open(configs.DBDriver, fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", configs.DBUser, configs.DBPassword, configs.DBHost, configs.DBPort, configs.DBName))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Initialize the repository
	orderRepo := database.NewOrderRepository(db)

	// Event dispatcher and order created event
	rabbitMQChannel := getRabbitMQChannel()
	eventDispatcher := events.NewEventDispatcher()
	eventDispatcher.Register("OrderCreated", &handler.OrderCreatedHandler{
		RabbitMQChannel: rabbitMQChannel,
	})

	orderCreatedEvent := events.NewOrderCreatedEvent()

	// Initialize use cases
	createOrderUseCase := usecase.NewCreateOrderUseCase(orderRepo, orderCreatedEvent, eventDispatcher)
	listOrdersUseCase := usecase.NewListOrdersUseCase(orderRepo)

	// Initialize Web Server and Handlers
	webOrderHandler := web.NewWebOrderHandler(eventDispatcher, orderRepo, orderCreatedEvent)

	// Use sync.WaitGroup to keep the main Goroutine alive
	var wg sync.WaitGroup
	wg.Add(3) // Add 3 for web, gRPC, and GraphQL servers

	// Start the Web Server (REST API)
	go func() {
		defer wg.Done()
		http.HandleFunc("/order", webOrderHandler.Create)
		http.HandleFunc("/orders", webOrderHandler.ListOrders)
		fmt.Println("Starting web server on port 8000")
		if err := http.ListenAndServe(":8000", nil); err != nil {
			fmt.Println("Error starting web server:", err)
		}
	}()

	// Start the gRPC Server
	go func() {
		defer wg.Done()
		fmt.Println("Starting gRPC server on port", configs.GRPCServerPort)
		lis, err := net.Listen("tcp", fmt.Sprintf(":%s", configs.GRPCServerPort))
		if err != nil {
			panic(err)
		}
		grpcServer := grpc.NewServer()
		createOrderService := service.NewOrderService(*createOrderUseCase)
		pb.RegisterOrderServiceServer(grpcServer, createOrderService)
		reflection.Register(grpcServer)
		if err := grpcServer.Serve(lis); err != nil {
			fmt.Println("Error starting gRPC server:", err)
		}
	}()

	// Start the GraphQL Server
	go func() {
		defer wg.Done()
		fmt.Println("Starting GraphQL server on port", configs.GraphQLServerPort)
		srv := graphql_handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{
			Resolvers: &graph.Resolver{
				CreateOrderUseCase: *createOrderUseCase,
				ListOrdersUseCase:  *listOrdersUseCase,
				OrderRepository:    orderRepo,
			},
		}))
		http.Handle("/", playground.Handler("GraphQL playground", "/query"))
		http.Handle("/query", srv)
		if err := http.ListenAndServe(":"+configs.GraphQLServerPort, nil); err != nil {
			fmt.Println("Error starting GraphQL server:", err)
		}
	}()

	// Block the main Goroutine by waiting for the servers to finish
	wg.Wait()
}

func getRabbitMQChannel() *amqp.Channel {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		panic(err)
	}
	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	return ch
}
