package order_status

import "time"

type OrderStatus struct {
	ID              int       `db:"id"`
	OrderId         int       `db:"order_id"`
	Status          string    `db:"status"`
	LastUpdatedDate time.Time `db:"last_update_date"`
}
