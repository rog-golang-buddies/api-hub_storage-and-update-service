package dto

import (
	"database/sql"
)

type ApiSpecDoc struct {
	ID          int            `db:"id"`
	Title       sql.NullString `db:"title"`
	Description sql.NullString `db:"description"`
	Type        sql.NullInt64  `db:"type"`
	Md5sum      sql.NullString `db:"md5sum"`
	CreatedAt   sql.NullTime   `db:"created_at"`
	UpdatedAt   sql.NullTime   `db:"updated_at"`
	FetchedAt   sql.NullTime   `db:"fetched_at"`
}
