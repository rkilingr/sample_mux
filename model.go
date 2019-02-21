package main

import (
	"database/sql"
	"fmt"
)

//Customer element struct
type Customer struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (c *Customer) getCustomer(db *sql.DB) error {
	query := "SELECT name,email,password FROM customers WHERE id=?"
	values := []interface{}{}
	values = append(values, c.ID)
	return db.QueryRow(query, values...).Scan(&c.Name, &c.Email, &c.Password)
}

func (c *Customer) createCustomer(db *sql.DB) error {
	query := "INSERT INTO customers(name,email,password) VALUES (?,?,?);"
	values := []interface{}{}
	values = append(values, c.Name, c.Email, c.Password)
	_, err := db.Exec(query, values...)

	if err != nil {
		return err
	}

	query = fmt.Sprintf("SELECT LAST_INSERT_ID()")
	err = db.QueryRow(query).Scan(&c.ID)

	if err != nil {
		return err
	}

	return nil
}

func (c *Customer) updateCustomer(db *sql.DB) error {
	query := "UPDATE customers SET name = ?, email = ?, password = ? where id=?"
	values := []interface{}{}
	values = append(values, c.Name, c.Email, c.Password, c.ID)
	_, err := db.Exec(query, values...)
	return err
}
func (c *Customer) deleteCustomer(db *sql.DB) error {
	query := "DELETE FROM customers WHERE id=?"
	values := []interface{}{}
	values = append(values, c.ID)
	_, err := db.Exec(query, values...)
	return err
}
func getCustomers(db *sql.DB, start, count int) ([]Customer, error) {
	query := "SELECT id,name,email,password FROM customers LIMIT ? OFFSET ?"
	values := []interface{}{}
	values = append(values, count, start)
	rows, err := db.Query(query, values...)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	customers := []Customer{}

	for rows.Next() {
		var c Customer
		if err := rows.Scan(&c.ID, &c.Name, &c.Email, &c.Password); err != nil {
			return nil, err
		}
		customers = append(customers, c)
	}
	return customers, nil
}
