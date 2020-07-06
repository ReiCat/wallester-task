package models

import (
	"database/sql"
	_ "github.com/lib/pq"
)

type Datastore interface {
	GetCustomers() ([]*Customer, error)
	GetCustomersByName(name, orderBy, order string, limit, skips int) ([]*Customer, error)
	DeleteCustomerById(id int) error
	AddCustomer(customer Customer) error
	GetCustomerById(id int) (Customer, error)
	UpdateCustomer(customer Customer) error
}

type DB struct {
	*sql.DB
}

func InitDB(dataSourceName string) (*DB, error) {
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return &DB{db}, nil
}
