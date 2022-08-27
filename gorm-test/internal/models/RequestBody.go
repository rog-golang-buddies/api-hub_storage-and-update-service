package dto

import (
	"database/sql"
)

type RequestBody struct {
	ID          int            `db:"id"`
	Description sql.NullString `db:"description"`
	Required    sql.NullBool   `db:"required"`
	ApiMethodID sql.NullInt64  `db:"api_method_id"`
}
