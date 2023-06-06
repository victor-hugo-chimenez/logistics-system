package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"io"
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

func getRoot(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	fmt.Printf("%s: got / request\n", ctx.Value(keyServerAddr))
	io.WriteString(w, "This is my website!\n")
}
func getHello(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	fmt.Printf("%s: got /hello request\n", ctx.Value(keyServerAddr))
	io.WriteString(w, "Hello, HTTP!\n")
}

func main() {
	log.Printf("Starting database\n")
	db := startDB()

	repository := delivery.NewRepository(db)
	service := delivery.NewDeliveryService(repository)
	controller := delivery.NewController(service)

	mux := http.NewServeMux()

	mux.HandleFunc("/delivery/:id", controller.FindById)
	mux.HandleFunc("/delivery")

	serverPort := "3000"

	ctx := context.Background()
	server := &http.Server{
		Addr:    serverPort,
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
