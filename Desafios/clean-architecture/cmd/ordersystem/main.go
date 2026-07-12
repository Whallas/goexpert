package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/graphql-go/handler"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/Whallas/goexpert/Desafios/clean-architecture/configs"
	"github.com/Whallas/goexpert/Desafios/clean-architecture/internal/infra/database"
	"github.com/Whallas/goexpert/Desafios/clean-architecture/internal/infra/graph"
	pb "github.com/Whallas/goexpert/Desafios/clean-architecture/internal/infra/grpc/pb"
	"github.com/Whallas/goexpert/Desafios/clean-architecture/internal/infra/grpc/service"
	web "github.com/Whallas/goexpert/Desafios/clean-architecture/internal/infra/web"
	"github.com/Whallas/goexpert/Desafios/clean-architecture/internal/infra/web/webserver"
	"github.com/Whallas/goexpert/Desafios/clean-architecture/internal/usecase"
)

func main() {
	cfg := configs.LoadConfig()

	db := openDBWithRetry(cfg)
	defer db.Close()

	runMigrations(cfg)

	orderRepo := database.NewOrderRepository(db)
	createUC := usecase.NewCreateOrderUseCase(orderRepo)
	listUC := usecase.NewListOrdersUseCase(orderRepo)

	go startWebServer(cfg, createUC, listUC)
	go startGRPCServer(cfg, createUC, listUC)
	startGraphQLServer(cfg, createUC, listUC)
}

func openDBWithRetry(cfg *configs.Conf) *sql.DB {
	const maxRetries = 30
	for i := 0; i < maxRetries; i++ {
		db, err := sql.Open(cfg.DBDriver, cfg.DBSource())
		if err == nil {
			if err = db.Ping(); err == nil {
				log.Println("database connected")
				return db
			}
			db.Close()
		}
		log.Printf("waiting for database... (%d/%d): %v", i+1, maxRetries, err)
		time.Sleep(2 * time.Second)
	}
	log.Fatal("database not ready after retries")
	return nil
}

func runMigrations(cfg *configs.Conf) {
	m, err := migrate.New("file://migrations", cfg.MigrateURL())
	if err != nil {
		log.Fatalf("migrate.New: %v", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("migrate.Up: %v", err)
	}
	log.Println("migrations applied")
}

func startWebServer(cfg *configs.Conf, createUC *usecase.CreateOrderUseCase, listUC *usecase.ListOrdersUseCase) {
	ws := webserver.NewWebServer(cfg.WebServerPort)
	orderHandler := web.NewOrderHandler(*createUC, *listUC)
	ws.AddHandler(http.MethodPost, "/order", orderHandler.Create)
	ws.AddHandler(http.MethodGet, "/order", orderHandler.List)
	log.Printf("REST server listening on :%s", cfg.WebServerPort)
	ws.Start()
}

func startGRPCServer(cfg *configs.Conf, createUC *usecase.CreateOrderUseCase, listUC *usecase.ListOrdersUseCase) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.GRPCServerPort))
	if err != nil {
		log.Fatalf("grpc listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	orderService := service.NewOrderService(*createUC, *listUC)
	pb.RegisterOrderServiceServer(grpcServer, orderService)
	reflection.Register(grpcServer)
	log.Printf("gRPC server listening on :%s", cfg.GRPCServerPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("grpc serve: %v", err)
	}
}

func startGraphQLServer(cfg *configs.Conf, createUC *usecase.CreateOrderUseCase, listUC *usecase.ListOrdersUseCase) {
	schema, err := graph.NewSchema(createUC, listUC)
	if err != nil {
		log.Fatalf("graphql schema: %v", err)
	}
	h := handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: true,
	})
	http.Handle("/graphql", h)
	log.Printf("GraphQL server listening on :%s/graphql", cfg.GraphQLServerPort)
	if err := http.ListenAndServe(":"+cfg.GraphQLServerPort, nil); err != nil {
		log.Fatalf("graphql serve: %v", err)
	}
}
