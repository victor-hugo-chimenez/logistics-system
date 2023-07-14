package order

import (
	"database/sql"
)

type Order struct {
	ID          int          `db:"id"`
	Amount      int          `db:"amount"`
	Description string       `db:"description"`
	UserId      int          `db:"user_id"`
	CreatedAt   sql.NullTime `db:"created_at"`
	UpdatedAt   sql.NullTime `db:"updated_at"`
}
