package delivery

type Delivery struct {
	ID       int    `db:"id"`
	OrderId  string `db:"order_id"`
	DriverId string `db:"driver_id"`
}
