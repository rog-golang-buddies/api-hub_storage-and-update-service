package dto

import (
	"database/sql"
)

type Schema struct {
	ID          int            `db:"id"`
	Key         sql.NullString `db:"key"`
	Type        sql.NullString `db:"type"`
	Description sql.NullString `db:"description"`
	ParentID    sql.NullInt64  `db:"parent_id"`
}
