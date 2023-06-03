package delivery

import (
	"context"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var Schema = `
	CREATE TABLE IF NOT EXISTS driver (
    id INT GENERATED ALWAYS AS IDENTITY,
	name TEXT,
    vehicle_model TEXT,
	PRIMARY KEY (id)
);
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
	if err := r.db.Get(&delivery, "SELECT * FROM delivery WHERE id=?", id); err != nil {
		return nil, err
	}
	return &delivery, nil
}
