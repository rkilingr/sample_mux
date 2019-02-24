package main

import "os"

//DbUser DataBase Username
var DbUser = os.Getenv("DB_USER")

//DbPassword Database password
var DbPassword = os.Getenv("DB_PASSWORD")

//DbHost DatabaseHost location
var DbHost = os.Getenv("DB_HOST")

//DbName DataBase Name
var DbName = "rest_api_example"

func main() {
	a := App{}
	a.Initialize(DbUser, DbPassword, DbHost, DbName)
	a.Run(":8080")
}
