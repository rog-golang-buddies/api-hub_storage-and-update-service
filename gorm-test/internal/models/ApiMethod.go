package dto

import (
	"database/sql"
)

type ApiMethod struct {
	ID           int            `db:"id"`
	Path         sql.NullString `db:"path"`
	Name         sql.NullString `db:"name"`
	Description  sql.NullString `db:"description"`
	Type         sql.NullString `db:"type"`
	ApiSpecDocID sql.NullInt64  `db:"api_spec_doc_id"`
	GroupID      sql.NullInt64  `db:"group_id"`
}
