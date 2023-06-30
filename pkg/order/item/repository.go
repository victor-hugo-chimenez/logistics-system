package order_item

import (
	"context"
	"github.com/jmoiron/sqlx"
	"log"
)

var Schema = `
	CREATE TABLE IF NOT EXISTS order_items (
    id INT GENERATED ALWAYS AS IDENTITY,
	order_id INT,
	item_name VARCHAR(255),
	quantity INT,    
	
	PRIMARY KEY (id),
	CONSTRAINT fk_order_id FOREIGN KEY (order_id) REFERENCES orders(id)	
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
func (r *Repository) FindItemByOrderId(ctx context.Context, id int) ([]OrderItem, error) {
	var orderItems []OrderItem
	if err := r.db.Select(&orderItems, "SELECT * FROM order_items WHERE order_id = $1", id); err != nil {
		return nil, err
	}
	return orderItems, nil
}

func (r *Repository) CreateOrderItem(ctx context.Context, item *OrderItem) error {
	tx := r.db.MustBegin()

	tx.MustExecContext(ctx, "INSERT INTO order_items (order_id, item_name, quantity) VALUES ($1, $2, $3)", item.OrderId, item.ItemName, item.Quantity)

	if err := tx.Commit(); err != nil {
		log.Fatalln(err)
	}

	return nil
}
