package driver

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
)

var Schema = `
	CREATE TABLE IF NOT EXISTS driver (
    id INT GENERATED ALWAYS AS IDENTITY,
	name VARCHAR(255),
    vehicle_model VARCHAR(255),
    vehicle_license_plate VARCHAR(7),
	license_number INT,
	    
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
	var limit int
	if err := r.db.Get(&limit, "SELECT COUNT(id) FROM driver"); err != nil {
		limit = 12
	}
	page := 4
	pageSize := limit / page

	if pageSize == 0 {
		pageSize = 1
	}

	fmt.Printf("Limit %d, Page %d, PageSize %d\n", limit, page, pageSize)

	for lastReadId := 0; lastReadId < limit; lastReadId += pageSize {
		var partialDriverResult []Driver
		if err := r.db.Select(&partialDriverResult, "SELECT * FROM driver WHERE id BETWEEN $1 AND $2 ORDER BY id DESC", lastReadId, lastReadId+pageSize); err != nil {
			return nil, err
		}
		fmt.Printf("Iteration %d", lastReadId)
		fmt.Println(partialDriverResult)
		drivers = append(drivers, partialDriverResult...)
	}

	return drivers, nil
}

func (r *Repository) CreateDriver(ctx context.Context, driver *Driver) error {
	tx := r.db.MustBegin()

	tx.MustExecContext(ctx, "INSERT INTO driver (name, vehicle_model, vehicle_license_plate, license_number) VALUES ($1, $2, $3, $4)", driver.Name, driver.VehicleModel, driver.VehicleLicensePlate, driver.LicenseNumber)

	if err := tx.Commit(); err != nil {
		log.Fatalln(err)
	}

	return nil
}

func (r *Repository) DeleteById(ctx context.Context, id int) error {
	tx := r.db.MustBegin()

	tx.MustExecContext(ctx, "DELETE FROM driver WHERE id = $1 ", id)

	if err := tx.Commit(); err != nil {
		log.Fatalln(err)
	}

	return nil
}
