package order_status

import "time"

type OrderStatus struct {
	Version   int       `db:"version"`
	OrderId   int       `db:"order_id"`
	Event     Event     `db:"event"`
	CreatedAt time.Time `db:"created_at"`
}

// Event TODO oq deveria ter dentro desse evento?
type Event struct {
	foo string
	bar string
}
