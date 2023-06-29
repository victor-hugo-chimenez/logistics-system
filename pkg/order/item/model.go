package order_item

type OrderItem struct {
	ID       int    `db:"id"`
	OrderId  int    `db:"order_id"`
	ItemName string `db:"item_name"`
	Quantity int    `db:"quantity"`
}
