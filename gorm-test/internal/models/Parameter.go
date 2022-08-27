package dto

import (
	"database/sql"
)

type Parameter struct {
	ID          int            `db:"id"`
	Name        sql.NullString `db:"name"`
	In          sql.NullString `db:"in"`
	Description sql.NullString `db:"description"`
	SchemaID    sql.NullInt64  `db:"schema_id"`
	Required    sql.NullBool   `db:"required"`
	ApiMethodID sql.NullInt64  `db:"api_method_id"`
}
