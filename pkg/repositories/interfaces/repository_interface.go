package repository_interfaces

import "github.com/jmoiron/sqlx"

type Repository interface {
	Create(*sqlx.DB) error
}
