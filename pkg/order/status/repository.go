package order_status

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
)

var Schema = `
	CREATE TABLE IF NOT EXISTS order_status_history (
	    version INT GENERATED ALWAYS AS IDENTITY,
	    order_id INT,
	    event jsonb,
	    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),
	    
	    PRIMARY KEY (order_id, version),
	    CONSTRAINT fk_order_id FOREIGN KEY (order_id) REFERENCES orders(id)                                          
	);

	CREATE TABLE IF NOT EXISTS order_status_checkpoint (
	    version INT,
	    order_id INT,
	    event jsonb,
	    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),
	    
	    UNIQUE (version),
	    CONSTRAINT fk_order_id FOREIGN KEY (order_id) REFERENCES orders(id)
	)
`

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		db,
	}
}

func (r *Repository) FindStatusByOrderId(ctx context.Context, id int) (*OrderStatus, error) {
	var orderStatus *OrderStatus
	if err := r.db.Select(&orderStatus, "SELECT * FROM order_status_checkpoint WHERE order_id = $1 ORDER BY version DESC LIMIT 1", id); err != nil {
		return nil, err
	}
	return orderStatus, nil
}

func (r *Repository) UpdateOrderStatusCheckpoint(ctx context.Context, status *OrderStatus) (*OrderStatus, error) {
	tx := r.db.MustBegin()

	var currentVersion *int
	err := tx.Get(&currentVersion, "SELECT MAX(version) from order_status_checkpoint where order_id=$1", status.OrderId)
	if err != nil {
		return nil, err
	}

	if _, err := tx.ExecContext(ctx, "INSERT INTO order_status_checkpoint (order_id, version, event) VALUES ($1, $2, $3)", status.OrderId, *currentVersion+1, status.Event); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to update order status %s", err)
	}

	var result *OrderStatus
	err = tx.Get(&result, "SELECT * from order_status_checkpoint where order_id=$1 order by version desc limit 1", status.OrderId)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *Repository) UpdateOrderStatusHistory(ctx context.Context, status *OrderStatus) error {
	tx := r.db.MustBegin()

	if _, err := tx.ExecContext(ctx, "INSERT INTO order_status_history (order_id, event) VALUES ($1, $2)", status.OrderId, status.Event); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to update order status %s", err)
	}

	return nil
}
