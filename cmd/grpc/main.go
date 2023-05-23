package main

import (
	"database/sql"
	"delivery-msg/config"
	"delivery-msg/internal/handlers"
	"delivery-msg/internal/repositories"
	"delivery-msg/internal/services"
	"delivery-msg/pb"
	"delivery-msg/pkg"
	"github.com/charmbracelet/log"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/nats-io/nats.go"
	"google.golang.org/grpc"
	"net"
)

func main() {

	cfg := config.ReadConfig()

	// connect DB
	databaseURL := cfg.DatabaseUrl
	if databaseURL == "" {
		log.Error("DATABASE_URL is not set")
		return
	}

	pgConn, err := sql.Open("pgx", databaseURL)
	if err != nil {
		log.Error(err)
		return
	}
	defer pgConn.Close()

	pgMigrate, err := migrate.New(pkg.MigrationPath, cfg.DatabaseUrl)
	if err != nil {
		log.Error(err)
		return
	}
	if err = pgMigrate.Up(); err != nil && err != migrate.ErrNoChange {
		log.Error(err)
		return
	}

	log.Info("gRPC API is running...")
	listener, err := net.Listen("tcp", pkg.GRPCHostPort)
	if err != nil {
		log.Error(err)
		return
	}

	natsConn, err := nats.Connect(pkg.NATSHostPost, nats.UserInfo(cfg.NATSUser, cfg.NATSPassword))
	if err != nil {
		log.Error(err)
		return
	}
	defer natsConn.Close()

	dbRepository := repositories.NewPostgreSQLRepository(pgConn)
	natsClient, err := services.NewNATSClient(natsConn)
	if err != nil {
		log.Error(err)
		return
	}

	deliveryHandler := handlers.NewDeliveryHandler(&dbRepository, &natsClient)

	s := grpc.NewServer()
	pb.RegisterDeliveryServiceServer(s, deliveryHandler)
	if err = s.Serve(listener); err != nil {
		log.Errorf("failed to serve %v", err)
		return
	}
}
