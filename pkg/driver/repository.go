package driver

import (
	"context"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
)

var Schema = `
	CREATE TABLE IF NOT EXISTS driver (
    id INT GENERATED ALWAYS AS IDENTITY,
	name TEXT,
    vehicle_model TEXT,
	PRIMARY KEY (id)
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

func (r *Repository) FindById(ctx context.Context, id int) (*Driver, error) {
	var driver Driver
	if err := r.db.Get(&driver, "SELECT * FROM driver WHERE id=$1", id); err != nil {
		return nil, err
	}
	return &driver, nil
}

func (r *Repository) FindAll(ctx context.Context) ([]Driver, error) {
	var drivers []Driver
	if err := r.db.Select(&drivers, "SELECT * FROM driver LIMIT 100"); err != nil {
		return nil, err
	}
	return drivers, nil
}

func (r *Repository) CreateDriver(ctx context.Context, driver *Driver) error {
	tx := r.db.MustBegin()

	tx.MustExecContext(ctx, "INSERT INTO driver (name, vehicle_model) VALUES ($1, $2)", driver.Name, driver.VehicleModel)

	if err := tx.Commit(); err != nil {
		log.Fatalln(err)
	}

	return nil
}
