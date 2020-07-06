package models

import (
	"time"
	"database/sql"
	"log"
)

type Customer struct {
	ID        int
	FirstName string
	LastName  string
	BirthDate string
	Gender    string
	Email     string
	Address   string
}

func parseDate(date string) string {
	dt, _ := time.Parse("2006-01-02T15:04:05Z", date)
	parsedDate := dt.Format("02.01.2006")
	return parsedDate
}

func (db *DB) GetCustomers() ([]*Customer, error) {
	sqlStatement := `SELECT id, first_name, last_name, 
    birth_date, gender, email, address FROM customers;`
	rows, err := db.Query(sqlStatement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	customers := make([]*Customer, 0)
	for rows.Next() {
		customer := new(Customer)
		err := rows.Scan(
			&customer.ID,
			&customer.FirstName,
			&customer.LastName,
			&customer.BirthDate,
			&customer.Gender,
			&customer.Email,
			&customer.Address)
		if err != nil {
			return nil, err
		}
		customer.BirthDate = parseDate(customer.BirthDate)
		customers = append(customers, customer)
	}
	if err = rows.Err(); err != nil {
		return customers, err
	}
	return customers, nil
}

func (db *DB) GetCustomersByName(name, orderBy, order string, limit, skips int) ([]*Customer, error) {
	sqlStatement := `SELECT id, first_name, last_name, 
	birth_date, gender, email, address FROM customers 
	WHERE first_name SIMILAR TO $1 OR last_name SIMILAR TO $1 
	ORDER BY ` + orderBy + ` ` + order + ` LIMIT $2 OFFSET $3;`
	rows, err := db.Query(sqlStatement, name, limit, skips)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	customers := make([]*Customer, 0)
	for rows.Next() {
		customer := new(Customer)
		err := rows.Scan(
			&customer.ID,
			&customer.FirstName,
			&customer.LastName,
			&customer.BirthDate,
			&customer.Gender,
			&customer.Email,
			&customer.Address)
		if err != nil {
			return nil, err
		}
		customer.BirthDate = parseDate(customer.BirthDate)
		customers = append(customers, customer)
	}
	if err = rows.Err(); err != nil {
		return customers, err
	}
	return customers, nil
}

func (db *DB) DeleteCustomerById(id int) error {
	sqlStatement := `DELETE FROM customers WHERE id=$1;`
	_, err := db.Exec(sqlStatement, id)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) GetCustomerById(id int) (Customer, error) {
	var customer Customer
	sqlStatement := `SELECT id, first_name, last_name, birth_date, gender, email, address FROM customers WHERE id=$1;`
	row := db.QueryRow(sqlStatement, id)
	err := row.Scan(
		&customer.ID,
		&customer.FirstName,
		&customer.LastName,
		&customer.BirthDate,
		&customer.Gender,
		&customer.Email,
		&customer.Address)
	if err != nil && err != sql.ErrNoRows {
		log.Println(err.Error())
		return customer, err
	}
	customer.BirthDate = parseDate(customer.BirthDate)
    return customer, nil
}

func (db *DB) AddCustomer(customer Customer) error {
	sqlStatement := `INSERT INTO customers (first_name, last_name,
		 birth_date, gender, email, address) VALUES ($1, $2, $3, $4, $5, $6);`
    _, err := db.Exec(sqlStatement,
        customer.FirstName,
        customer.LastName,
        customer.BirthDate,
        customer.Gender,
        customer.Email,
        customer.Address)
    if err != nil {
        return err
    }
    return nil
}

func (db *DB) UpdateCustomer(customer Customer) error {
	sqlStatement := `UPDATE customers SET first_name=$2, last_name=$3,
	 birth_date=$4, gender=$5, email=$6, address=$7 WHERE id=$1;`

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	_, err = tx.Exec(sqlStatement,
		customer.ID,
        customer.FirstName,
        customer.LastName,
        customer.BirthDate,
        customer.Gender,
        customer.Email,
        customer.Address)
    if err != nil {
        return err
	}
	
	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}
    return nil
}
