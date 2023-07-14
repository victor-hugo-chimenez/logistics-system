package products

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
)

var Schema = `
	CREATE TABLE IF NOT EXISTS products (    
    id INT GENERATED ALWAYS AS IDENTITY,
 	price INT,
 	name VARCHAR(255),
 	description VARCHAR(255),
	discount_percent INT,    
    order_item_id INT,
    PRIMARY KEY (id),
	CONSTRAINT fk_order_item_id FOREIGN KEY (order_item_id) REFERENCES order_items(id)
);`

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		db,
	}
}

func (r *Repository) FindById(ctx context.Context, id int) (*Product, error) {
	var delivery Product
	if err := r.db.Get(&delivery, "SELECT * FROM delivery WHERE id=$1", id); err != nil {
		return nil, err
	}
	return &delivery, nil
}

func (r *Repository) FindAll(ctx context.Context) ([]Product, error) {
	var products []Product
	var limit int
	if err := r.db.Get(&limit, "SELECT COUNT(id) FROM products"); err != nil {
		limit = 12
	}
	page := 4
	pageSize := limit / page

	if pageSize == 0 {
		pageSize = 1
	}

	fmt.Printf("Limit %d, Page %d, PageSize %d\n", limit, page, pageSize)

	for lastReadId := 0; lastReadId < limit; lastReadId += pageSize {
		var partialProductResult []Product
		if err := r.db.Select(&partialProductResult, "SELECT * FROM products WHERE id BETWEEN $1 AND $2 ORDER BY id DESC", lastReadId, lastReadId+pageSize); err != nil {
			return nil, err
		}
		fmt.Printf("Iteration %d", lastReadId)
		fmt.Printf("Limit %d, Page %d, PageSize %d\n", limit, page, pageSize)
		fmt.Println(partialProductResult)
		products = append(products, partialProductResult...)
	}

	return products, nil
}

func (r *Repository) CreateProduct(ctx context.Context, product *Product) error {
	tx := r.db.MustBegin()

	tx.MustExecContext(ctx, "INSERT INTO products (price, name, description, discount_percent, order_item_id) VALUES ($1, $2, $3, $4, $5)", product.Price, product.Name, product.Description, product.DiscountPercent, product.OrderItemId)

	if err := tx.Commit(); err != nil {
		log.Fatalln(err)
	}

	return nil
}

func (r *Repository) UpdateProduct(ctx context.Context, product *Product) (*Product, error) {
	tx := r.db.MustBegin()

	tx.MustExecContext(ctx, "UPDATE products SET price = $2, name = $3, description = $4, discount_percent = $5, order_item_id = $6 WHERE id = $1", product.ID, product.Price, product.Name, product.Description, product.DiscountPercent, product.OrderItemId)

	if err := tx.Commit(); err != nil {
		log.Fatalln(err)
	}

	var updatedProduct Product
	if err := r.db.Get(&updatedProduct, "SELECT * FROM products WHERE id=$1", product.ID); err != nil {
		return nil, err
	}

	return &updatedProduct, nil
}

func (r *Repository) DeleteById(ctx context.Context, id int) error {
	tx := r.db.MustBegin()

	fmt.Println(id)
	tx.MustExecContext(ctx, "DELETE FROM products WHERE id = $1 ", id)

	if err := tx.Commit(); err != nil {
		log.Fatalln(err)
	}

	return nil
}
