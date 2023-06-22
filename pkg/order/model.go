package order

type Order struct {
	ID          int    `db:"id"`
	Amount      int    `db:"amount"`
	Description string `db:"driver_id"`
	DeliveryId  int    `db:"delivery_id"`
}
