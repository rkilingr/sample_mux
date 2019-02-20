package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
)

var a App

//EndPoint Test
func TestEmptyTable(t *testing.T) {
	clearTable()
	req, _ := http.NewRequest("GET", "/customers", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
	if body := response.Body.String(); body != "[]" {
		t.Errorf("Expected empty array. Received %s\n", body)
	}
}

func TestCreateCustomer(t *testing.T) {
	clearTable()
	payload := []byte(`{"name":"test cust", "email":"test@test.com","password":"pass"}`)
	req, _ := http.NewRequest("POST", "/customer", bytes.NewBuffer(payload))
	response := executeRequest(req)
	checkResponseCode(t, http.StatusCreated, response.Code)
	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["name"] != "test cust" {
		t.Errorf("Expected name to be %s, Got %s", "test_cust", m["name"])
	}
	if m["email"] != "test@test.com" {
		t.Errorf("Expected email to be %s, Got %s", "test@test.com", m["email"])
	}
	if m["password"] != "pass" {
		t.Errorf("Expected password to be %s, Got %s", "pass", m["password"])
	}
	if m["id"] != 1.0 {
		t.Errorf("Expected id to be %f, Got %f", 1.0, m["id"])
	}

}

func TestGetCustomer(t *testing.T) {
	clearTable()
	addCustomers(1)

	req, _ := http.NewRequest("GET", "/customer/1", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

}

func addCustomers(count int) {
	if count < 1 {
		count = 1
	}
	for i := 1; i <= count; i++ {
		statement := fmt.Sprintf("INSERT INTO customers(name,email,password) VALUES(%s,%s,%s)", "cust "+strconv.Itoa(i), "test"+strconv.Itoa(i)+"@test.com", "pass"+strconv.Itoa(i))
		a.DB.Exec(statement)
	}
}

func TestUpdateCustomer(t *testing.T) {
	clearTable()
	addCustomers(1)

	req, _ := http.NewRequest("GET", "/customer/1", nil)
	response := executeRequest(req)
	var originalCustomer map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &originalCustomer)

	payload := []byte(`{"name":"test cust-updated name", "email":"testupdated@test.com", "password":"passupdated"}`)

	req, _ = http.NewRequest("PUT", "/customer/1", bytes.NewBuffer(payload))
	response = executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["id"] != originalCustomer["id"] {
		t.Errorf("Expected the id to remain the same(%f), but Got %f\n", m["id"], originalCustomer["id"])
	}
	if m["name"] == originalCustomer["name"] {
		t.Errorf("Expected the name to be updated to %s,  but Got %s\n", m["name"], originalCustomer["name"])
	}
	if m["email"] == originalCustomer["email"] {
		t.Errorf("Expected the email to be updated to %s,  but Got %s\n", m["email"], originalCustomer["email"])
	}
	if m["password"] == originalCustomer["password"] {
		t.Errorf("Expected the password to be updated to %s,  but Got %s\n", m["password"], originalCustomer["password"])
	}

}

func TestDeleteCustomer(t *testing.T) {
	clearTable()
	addCustomers(1)

	req, _ := http.NewRequest("GET", "/customer/1", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("DELETE", "/customer/1", nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("GET", "/customer/1", nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)

}

func TestGetNonExistentCustomer(t *testing.T) {
	clearTable()
	req, _ := http.NewRequest("GET", "/customer/45", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "Customer not found" {
		t.Errorf("Expected the 'error' key of the response to be set to 'Customer not found'. Got '%s'", m["error"])
	}
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)
	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func TestMain(m *testing.M) {
	a = App{}
	a.Initialize("ravi", "hello123", "rest_api_example")
	ensureTableExists()
	code := m.Run()
	clearTable()
	os.Exit(code)
}

func ensureTableExists() {
	if _, err := a.DB.Exec(tableCreationQuery); err != nil {
		log.Println(err)
	}
}
func clearTable() {
	a.DB.Exec("DELETE FROM customers")
	a.DB.Exec("ALTER TABLE customers AUTO_INCREMENT=1")
}

const tableCreationQuery = `CREATE TABLE customers (
  id INT AUTO_INCREMENT PRIMARY KEY,
  name VARCHAR(50),
  email VARCHAR(50),
  password VARCHAR(50)
  )`
