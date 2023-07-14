package order_status

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
)

var Schema = `
	CREATE TABLE IF NOT EXISTS order_status_history (
	    version INT GENERATED ALWAYS AS IDENTITY,
	    order_id INT REFERENCES orders(id),
	    event jsonb,
	    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),
	    
	    PRIMARY KEY (order_id, version)
	);

	CREATE TABLE IF NOT EXISTS order_status_checkpoint (
	    version INT,
	    order_id INT,
	    event jsonb,
	    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),
	    
	    PRIMARY KEY (order_id, version),
	    CONSTRAINT fk_order_status_version FOREIGN KEY (version) REFERENCES order_status_history(version),
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

func (r *Repository) FindStatusByOrderId(ctx context.Context, id int) ([]OrderStatus, error) {
	var orderStatus []OrderStatus
	if err := r.db.Select(&orderStatus, "SELECT * FROM order_status WHERE order_id = $1", id); err != nil {
		return nil, err
	}
	return orderStatus, nil
}

func (r *Repository) UpdateOrderStatus(ctx context.Context, status *OrderStatus) error {
	tx := r.db.MustBegin()

	if _, err := tx.ExecContext(ctx, "INSERT INTO order_status_history (order_id, status) VALUES ($1, $2)", status.OrderId, status.Status); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to update order status %s", err)
	}

	return nil
}

// TODO Append no log: Todo comando recebido de update order deve guardar o evento
// event -> command
// Checkpoint com CronJob ou com "múltiplo de 10"
