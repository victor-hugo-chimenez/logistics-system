package order

import "time"

type Order struct {
	ID          int       `db:"id"`
	Amount      int       `db:"amount"`
	Description string    `db:"description"`
	UserId      int       `db:"user_id"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}
