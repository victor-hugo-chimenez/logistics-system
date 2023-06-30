package delivery

import "time"

type Delivery struct {
	ID        int       `db:"id"`
	OrderId   string    `db:"order_id"`
	Status    string    `db:"status"`
	DriverId  string    `db:"driver_id"`
	StartTime time.Time `db:"start_time"`
	EndTime   time.Time `db:"end_time"`
}
