package dto

import (
	"database/sql"
)

type Server struct {
	ID          int            `db:"id"`
	URL         sql.NullString `db:"url"`
	Description sql.NullString `db:"description"`
	ApiMethodID sql.NullInt64  `db:"api_method_id"`
}
