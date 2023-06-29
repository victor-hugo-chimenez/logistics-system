package user

import "github.com/jmoiron/sqlx"

var Schema = `
	CREATE TABLE IF NOT EXISTS users (    
    id INT GENERATED ALWAYS AS IDENTITY,
 	name VARCHAR(255),
 	email VARCHAR(255),
	address VARCHAR(255),    
	phone_number VARCHAR(255),    
    order_item_id INT,
    PRIMARY KEY (id)
);`

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		db,
	}
}
