package dto

import (
	"database/sql"
)

type ExternalDoc struct {
	ID          int            `db:"id"`
	Description sql.NullString `db:"description"`
	URL         sql.NullString `db:"url"`
	ApiMethodID sql.NullInt64  `db:"api_method_id"`
}
