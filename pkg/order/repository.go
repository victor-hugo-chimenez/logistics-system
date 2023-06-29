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
	user_id INT,
	created_at TIMESTAMP NOT NULL DEFAULT NOW(),    
	updated_at TIMESTAMP NOT NULL DEFAULT NOW(),   
	    
    PRIMARY KEY (id),
	CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(id)
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
	if err := r.db.Get(&order, "SELECT * FROM orders WHERE id=$1", id); err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *Repository) FindAll(ctx context.Context) ([]Order, error) {
	var order []Order
	var limit int
	if err := r.db.Get(&limit, "SELECT COUNT(id) FROM orders"); err != nil {
		limit = 12
	}
	page := 4
	pageSize := limit / page

	if pageSize == 0 {
		pageSize = 1
	}

	fmt.Printf("Limit %d, Page %d, PageSize %d\n", limit, page, pageSize)

	for lastReadId := 0; lastReadId < limit; lastReadId += pageSize {
		var partialOrderResult []Order
		if err := r.db.Select(&partialOrderResult, "SELECT * FROM orders WHERE id BETWEEN $1 AND $2 ORDER BY id DESC", lastReadId, lastReadId+pageSize); err != nil {
			fmt.Printf("Error getting orders %s", err)
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

	tx.MustExecContext(ctx, "INSERT INTO orders (amount, description, user_id, created_at, updated_at) VALUES ($1, $2, $3)", order.Amount, order.Description, order.UserId, order.CreatedAt, order.UpdatedAt)

	if err := tx.Commit(); err != nil {
		log.Fatalln(err)
	}

	return nil
}

func (r *Repository) UpdateOrder(ctx context.Context, order *Order) (*Order, error) {
	tx := r.db.MustBegin()

	tx.MustExecContext(ctx, "UPDATE orders SET amount = $2, description = $3, user_id = $4, updated_at = NOW() WHERE id = $1", order.Amount, order.Description, order.UserId)

	if err := tx.Commit(); err != nil {
		log.Fatalln(err)
	}

	var updatedOrder Order
	if err := r.db.Get(&updatedOrder, "SELECT * FROM orders WHERE id=$1", order.ID); err != nil {
		return nil, err
	}

	return &updatedOrder, nil
}

func (r *Repository) DeleteById(ctx context.Context, id int) error {
	tx := r.db.MustBegin()

	fmt.Println(id)
	tx.MustExecContext(ctx, "DELETE FROM orders WHERE id = $1 ", id)

	if err := tx.Commit(); err != nil {
		log.Fatalln(err)
	}

	return nil
}
