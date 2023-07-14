package products

type Product struct {
	ID              int    `db:"id"`
	Price           int    `db:"price"`
	Name            string `db:"name"`
	Description     string `db:"description"`
	DiscountPercent int    `db:"discount_percent"`
	OrderItemId     int    `db:"order_item_id"`
}
