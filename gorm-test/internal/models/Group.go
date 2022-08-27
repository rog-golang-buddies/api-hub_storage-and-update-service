package dto

import (
	"database/sql"
)

type Group struct {
	ID           int            `db:"id"`
	Name         sql.NullString `db:"name"`
	Description  sql.NullString `db:"description"`
	ApiSpecDocID sql.NullInt64  `db:"api_spec_doc_id"`
}
