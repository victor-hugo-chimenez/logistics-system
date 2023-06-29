package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
	"logistics_system/pkg/delivery"
	"logistics_system/pkg/driver"
	"logistics_system/pkg/order"
	"logistics_system/pkg/order/item"
	order_status "logistics_system/pkg/order/status"
	"logistics_system/pkg/products"
	"logistics_system/pkg/user"
	"net"
	"net/http"
)

const (
	dbHost        = "localhost"
	dbPort        = 5432
	dbUser        = "postgres"
	dbPassword    = "postgres"
	dbName        = "postgres"
	keyServerAddr = "serverAddr"
)

func startDB() *sqlx.DB {
	pgInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	// this Pings the database trying to connect
	// use sqlx.Open() for sql.Open() semantics
	db, err := sqlx.Connect("postgres", pgInfo)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Exec user")
	db.MustExec(user.Schema)
	fmt.Println("Exec order")
	db.MustExec(order.Schema)
	fmt.Println("Exec item")
	db.MustExec(order_item.Schema)
	fmt.Println("Exec status")
	db.MustExec(order_status.Schema)
	fmt.Println("Exec prduct")
	db.MustExec(products.Schema)
	fmt.Println("Exec driver")
	db.MustExec(driver.Schema)
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
