package order_status

import (
	"context"
	"github.com/jmoiron/sqlx"
	"log"
)

var Schema = `
	CREATE TABLE IF NOT EXISTS order_status (
    id INT GENERATED ALWAYS AS IDENTITY,
	order_id INT,
	status VARCHAR(255),
	last_update_date TIMESTAMP NOT NULL DEFAULT NOW(),    
	
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

func (r *Repository) FindStatusByOrderId(ctx context.Context, id int) ([]OrderStatus, error) {
	var orderStatus []OrderStatus
	if err := r.db.Select(&orderStatus, "SELECT * FROM order_status WHERE order_id = $1", id); err != nil {
		return nil, err
	}
	return orderStatus, nil
}

func (r *Repository) CreateOrderStatus(ctx context.Context, status *OrderStatus) error {
	tx := r.db.MustBegin()

	tx.MustExecContext(ctx, "INSERT INTO order_status (order_id, status) VALUES ($1, $2)", status.OrderId, status.Status)

	if err := tx.Commit(); err != nil {
		log.Fatalln(err)
	}

	return nil
}
