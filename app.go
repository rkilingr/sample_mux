package main

// Create a new user in response to a valid POST request at /user,
// Update a user in response to a valid PUT request at /user/{id},
// Delete a user in response to a valid DELETE request at /user/{id},
// Fetch a user in response to a valid GET request at /user/{id}, and
// Fetch a list of users in response to a valid GET request at /users.
import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"github.com/jinzhu/gorm"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

//App Struct for containing the router and DB for the project
type App struct {
	Router *mux.Router
	DB     *gorm.DB
}

//Initialize Initializes the DB for the project t be ready
func (a *App) Initialize(user, password, dbhost, dbname string) {
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8&parseTime=True&loc=Local", user, password, dbhost, dbname)
	fmt.Println(connectionString)
	var err error
	a.DB, err = gorm.Open("mysql", connectionString)
	if a.DB.HasTable(&Customer{}){
		a.DB.AutoMigrate(&Customer{})
	}else{
		a.DB.CreateTable(&Customer{})
	}
	if err != nil {
		log.Fatal(err)
	}
	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

func (a *App) getCustomer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid User ID")
		return
	}
	customer := Customer{ID: uint(id)}
	if err := customer.getCustomer(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Customer not found")
		default:
			respondWithError(w, http.StatusBadRequest, "Bad request")
		}
		return
	}
	respondWithJSON(w, http.StatusOK, customer)
}

func respondWithError(w http.ResponseWriter, status int, error string) {
	respondWithJSON(w, status, map[string]string{"error": error})
}

func respondWithJSON(w http.ResponseWriter, status int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(response)
}

func (a *App) getCustomers(w http.ResponseWriter, r *http.Request) {
	count, _ := strconv.Atoi(r.FormValue("count"))
	start, _ := strconv.Atoi(r.FormValue("start"))
	if count > 10 || count < 1 {
		count = 10
	}
	if start < 0 {
		start = 0
	}
	customers, err := getCustomers(a.DB, start, count)

	if err != nil {
		respondWithError(w, http.StatusNotFound, "Customers not found")
		return
	}
	respondWithJSON(w, http.StatusOK, customers)
}

func (a *App) createCustomer(w http.ResponseWriter, r *http.Request) {
	var c Customer
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&c); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Request payload")
		return
	}
	defer r.Body.Close()

	if err := c.createCustomer(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Server Error")
		return
	}
	respondWithJSON(w, http.StatusCreated, c)
}

func (a *App) updateCustomer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid customer ID")
		return
	}
	var c Customer
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&c); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()
	c.ID = uint(id)

	if err := c.updateCustomer(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Server Error")
		return
	}

	respondWithJSON(w, http.StatusOK, c)
}

func (a *App) deleteCustomer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Customer ID")
		return
	}
	c := Customer{ID: uint(id)}

	if err := c.deleteCustomer(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Server error")
	}
	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

//initializeRoutes Initializes all the routes
func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/customer", a.getCustomers).Methods("GET")
	a.Router.HandleFunc("/customer/{id:[0-9]+}", a.getCustomer).Methods("GET")
	a.Router.HandleFunc("/customer/{id:[0-9]+}", a.updateCustomer).Methods("PUT")
	a.Router.HandleFunc("/customer/{id:[0-9]+}", a.deleteCustomer).Methods("DELETE")
	a.Router.HandleFunc("/customer", a.createCustomer).Methods("POST")

}

//Run Runs the web backend
func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}
