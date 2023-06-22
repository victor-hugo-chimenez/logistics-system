package products

import "github.com/jmoiron/sqlx"

var Schema = `
	CREATE TABLE IF NOT EXISTS products (    
    id INT GENERATED ALWAYS AS IDENTITY,
 	price INT,
 	name VARCHAR(255),
 	description VARCHAR(255),
	discount_percent INT,    
    order_item_id INT,
    PRIMARY KEY (id),
	CONSTRAINT fk_order_item_id FOREIGN KEY (order_item) REFERENCES order_item(id)
);`

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		db,
	}
}
