package driver

type Driver struct {
	ID           int `db:"id"`
	Name         string
	VehicleModel string `db:"vehicle_model"`
}
