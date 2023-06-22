package delivery

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
)

var Schema = `
	CREATE TABLE IF NOT EXISTS delivery (    
    id INT GENERATED ALWAYS AS IDENTITY,
	order_id INT,
    driver_id INT,
    PRIMARY KEY (id),
	CONSTRAINT fk_driver_id FOREIGN KEY (driver_id) REFERENCES driver(id)
);

`

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		db,
	}
}

func (r *Repository) FindById(ctx context.Context, id int) (*Delivery, error) {
	var delivery Delivery
	if err := r.db.Get(&delivery, "SELECT * FROM delivery WHERE id=$1", id); err != nil {
		return nil, err
	}
	return &delivery, nil
}

func (r *Repository) FindAll(ctx context.Context) ([]Delivery, error) {
	var delivery []Delivery
	var limit int
	if err := r.db.Get(&limit, "SELECT COUNT(id) FROM delivery"); err != nil {
		limit = 12
	}
	page := 4
	pageSize := limit / page
	fmt.Printf("Limit %d, Page %d, PageSize %d\n", limit, page, pageSize)

	for lastReadId := 0; lastReadId < limit; lastReadId += pageSize {
		var partialDeliveryResult []Delivery
		if err := r.db.Select(&partialDeliveryResult, "SELECT * FROM delivery WHERE id BETWEEN $1 AND $2 ORDER BY id DESC", lastReadId, lastReadId+pageSize); err != nil {
			return nil, err
		}
		fmt.Printf("Iteration %d", lastReadId)
		fmt.Println(partialDeliveryResult)
		delivery = append(delivery, partialDeliveryResult...)
	}

	return delivery, nil
}

func (r *Repository) CreateDelivery(ctx context.Context, delivery *Delivery) error {
	tx := r.db.MustBegin()

	tx.MustExecContext(ctx, "INSERT INTO delivery (order_id, driver_id) VALUES ($1, $2)", delivery.OrderId, delivery.DriverId)

	if err := tx.Commit(); err != nil {
		log.Fatalln(err)
	}

	return nil
}

func (r *Repository) UpdateDelivery(ctx context.Context, delivery *Delivery) (*Delivery, error) {
	tx := r.db.MustBegin()

	tx.MustExecContext(ctx, "UPDATE delivery SET order_id = $2, driver_id = $3 WHERE id = $1", delivery.ID, delivery.OrderId, delivery.DriverId)

	if err := tx.Commit(); err != nil {
		log.Fatalln(err)
	}

	var updatedDelivery Delivery
	if err := r.db.Get(&updatedDelivery, "SELECT * FROM delivery WHERE id=$1", delivery.ID); err != nil {
		return nil, err
	}

	return &updatedDelivery, nil
}

func (r *Repository) DeleteById(ctx context.Context, id int) error {
	tx := r.db.MustBegin()

	fmt.Println(id)
	tx.MustExecContext(ctx, "DELETE FROM delivery WHERE id = $1 ", id)

	if err := tx.Commit(); err != nil {
		log.Fatalln(err)
	}

	return nil
}
