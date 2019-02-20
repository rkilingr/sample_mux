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
	query := fmt.Sprintf("SELECT name,email,password FROM customers WHERE id=%d", c.ID)
	return db.QueryRow(query).Scan(&c.Name, &c.Email, &c.Password)
}

func (c *Customer) createCustomer(db *sql.DB) error {
	query := fmt.Sprintf("INSERT INTO customers(name,email,password) VALUES (%s,%s,%s)", c.Name, c.Email, c.Password)
	_, err := db.Exec(query)

	if err != nil {
		return err
	}

	query = fmt.Sprintf("SELECT LAST_INSERTED_ID()")
	err = db.QueryRow(query).Scan(&c.ID)

	if err != nil {
		return err
	}

	return nil
}

func (c *Customer) updateCustomer(db *sql.DB) error {
	query := fmt.Sprintf("UPDATE customers SET name = %s, email = %s, password = %s where id=%d", c.Name, c.Email, c.Password, c.ID)
	_, err := db.Exec(query)
	return err
}
func (c *Customer) deleteCustomer(db *sql.DB) error {
	query := fmt.Sprintf("DELETE customers where id=%d", c.ID)
	_, err := db.Exec(query)
	return err
}
func getCustomers(db *sql.DB, start, count int) ([]Customer, error) {
	query := fmt.Sprintf("SELECT id,name,email,password FROM customers LIMIT %d OFFSET %d", count, start)
	rows, err := db.Query(query)

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
