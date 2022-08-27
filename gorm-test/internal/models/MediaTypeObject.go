package dto

import (
	"database/sql"
)

type MediaTypeObject struct {
	ID            int           `db:"id"`
	RequestBodyID sql.NullInt64 `db:"request_body_id"`
	SchemaID      sql.NullInt64 `db:"schema_id"`
}
