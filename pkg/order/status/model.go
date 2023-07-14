package order_status

import (
	"database/sql"
)

type OrderStatus struct {
	ID              int          `db:"id"`
	OrderId         int          `db:"order_id"`
	Status          string       `db:"status"`
	LastUpdatedDate sql.NullTime `db:"last_update_date"`
}
