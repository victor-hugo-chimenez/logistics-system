package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
	"logistics_system/pkg/delivery"
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

	repository := delivery.NewRepository(db)
	service := delivery.NewDeliveryService(repository)
	controller := delivery.NewController(service)

	mux := http.NewServeMux()

	mux.HandleFunc("/delivery", controller.HandleRequest)

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

	err := server.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error listening for server: %s\n", err)
	}
}
