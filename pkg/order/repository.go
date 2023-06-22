package order

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
)

var Schema = `
	CREATE TABLE IF NOT EXISTS orders (    
    id INT GENERATED ALWAYS AS IDENTITY,
 	amount INT,
 	description VARCHAR(255),
    delivery_id INT,
    PRIMARY KEY (id),
	CONSTRAINT fk_delivery_id FOREIGN KEY (delivery_id) REFERENCES delivery(id)
);`

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		db,
	}
}

func (r *Repository) FindById(ctx context.Context, id int) (*Order, error) {
	var order Order
	if err := r.db.Get(&order, "SELECT * FROM order WHERE id=$1", id); err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *Repository) FindAll(ctx context.Context) ([]Order, error) {
	var order []Order
	var limit int
	if err := r.db.Get(&limit, "SELECT COUNT(id) FROM order"); err != nil {
		limit = 12
	}
	page := 4
	pageSize := limit / page
	fmt.Printf("Limit %d, Page %d, PageSize %d\n", limit, page, pageSize)

	for lastReadId := 0; lastReadId < limit; lastReadId += pageSize {
		var partialOrderResult []Order
		if err := r.db.Select(&partialOrderResult, "SELECT * FROM order WHERE id BETWEEN $1 AND $2 ORDER BY id DESC", lastReadId, lastReadId+pageSize); err != nil {
			return nil, err
		}
		fmt.Printf("Iteration %d", lastReadId)
		fmt.Println(partialOrderResult)
		order = append(order, partialOrderResult...)
	}

	return order, nil
}

func (r *Repository) CreateOrder(ctx context.Context, order *Order) error {
	tx := r.db.MustBegin()

	tx.MustExecContext(ctx, "INSERT INTO order (amount, description, delivery_id) VALUES ($1, $2, $3)", order.Amount, order.Description, order.DeliveryId)

	if err := tx.Commit(); err != nil {
		log.Fatalln(err)
	}

	return nil
}

func (r *Repository) UpdateOrder(ctx context.Context, order *Order) (*Order, error) {
	tx := r.db.MustBegin()

	tx.MustExecContext(ctx, "UPDATE order SET amount = $2, description = $3, delivery_id = $4 WHERE id = $1", order.ID, order.Amount, order.Description, order.DeliveryId)

	if err := tx.Commit(); err != nil {
		log.Fatalln(err)
	}

	var updatedOrder Order
	if err := r.db.Get(&updatedOrder, "SELECT * FROM order WHERE id=$1", order.ID); err != nil {
		return nil, err
	}

	return &updatedOrder, nil
}

func (r *Repository) DeleteById(ctx context.Context, id int) error {
	tx := r.db.MustBegin()

	fmt.Println(id)
	tx.MustExecContext(ctx, "DELETE FROM order WHERE id = $1 ", id)

	if err := tx.Commit(); err != nil {
		log.Fatalln(err)
	}

	return nil
}
