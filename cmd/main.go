package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
	"logistics_system/pkg/delivery"
	"logistics_system/pkg/driver"
	"net"
	"net/http"
)

const (
	host          = "localhost"
	port          = 5432
	user          = "postgres"
	password      = "postgres"
	dbname        = "postgres"
	keyServerAddr = "serverAddr"
)

func startDB() *sqlx.DB {
	pgInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// this Pings the database trying to connect
	// use sqlx.Open() for sql.Open() semantics
	db, err := sqlx.Connect("postgres", pgInfo)
	if err != nil {
		log.Fatalln(err)
	}

	db.MustExec(delivery.Schema)
	return db
}

func main() {
	log.Printf("Starting database\n")
	db := startDB()

	deliveryRepository := delivery.NewRepository(db)
	deliveryService := delivery.NewDeliveryService(deliveryRepository)
	deliveryController := delivery.NewController(deliveryService)

	driverRepository := driver.NewRepository(db)
	driverService := driver.NewDriverService(driverRepository)
	driverController := driver.NewController(driverService)

	mux := http.NewServeMux()

	mux.HandleFunc("/delivery", deliveryController.NewRouter())
	mux.HandleFunc("/driver", driverController.HandleDriverRequest)

	serverAddress := "127.0.0.1:3000"

	ctx := context.Background()
	server := &http.Server{
		Addr:    serverAddress,
		Handler: mux,
		BaseContext: func(l net.Listener) context.Context {
			ctx = context.WithValue(ctx, keyServerAddr, l.Addr().String())
			return ctx
		},
	}

	log.Printf("Server listening on %s", serverAddress)
	err := server.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error listening for server: %s\n", err)
	}
}
