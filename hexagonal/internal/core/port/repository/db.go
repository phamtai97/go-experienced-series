package repository

import (
	"database/sql"
	"io"
)

type Database interface {
	io.Closer
	GetDB() *sql.DB
}
