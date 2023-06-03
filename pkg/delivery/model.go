package delivery

type Driver struct {
	ID           int `db:"id"`
	Name         string
	VehicleModel string `db:"vehicle_model"`
}

type Delivery struct {
	ID       int    `db:"id"`
	OrderId  string `db:"order_id"`
	DriverId string `db:"driver_id"`
}
