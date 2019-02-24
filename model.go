package main

import (
	"database/sql"
	"time"

	"github.com/jinzhu/gorm"
)

//Customer element struct
type Customer struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

func (c *Customer) getCustomer(db *gorm.DB) error {
	if db.First(&c).RecordNotFound() {
		return sql.ErrNoRows
	}
	query := "SELECT * FROM customers where id=?"
	values := []interface{}{}
	values = append(values, c.ID)
	err := db.Raw(query, values...).Scan(&c)
	return err.Error
}

func (c *Customer) createCustomer(db *gorm.DB) error {
	//query := "INSERT INTO customers(name,email,password) VALUES (?,?,?);"
	err := db.Create(&c)

	if err != nil {
		return err.Error
	}

	//query = fmt.Sprintf("SELECT LAST_INSERT_ID()")
	//err = db.Query(query).Scan(&c.ID)
	err = db.Last(&Customer{}).Scan(&c)

	if err != nil {
		return err.Error
	}

	return nil
}

func (c *Customer) updateCustomer(db *gorm.DB) error {
	newC := Customer{ID: c.ID}
	db.First(&newC).Scan(&newC)
	newC.Name = c.Name
	newC.Email = c.Email
	newC.Password = c.Password
	err := db.Save(&newC)
	return err.Error
}
func (c *Customer) deleteCustomer(db *gorm.DB) error {
	err := db.Delete(&c)
	return err.Error
}
func getCustomers(db *gorm.DB, start, count int) ([]Customer, error) {
	query := "SELECT id,name,email,password FROM customers LIMIT ? OFFSET ?"
	values := []interface{}{}
	values = append(values, count, start)

	customers := []Customer{}
	err := db.Raw(query, values...).Scan(&customers)

	if err.Error != nil {
		return nil, err.Error
	}

	return customers, nil
}
