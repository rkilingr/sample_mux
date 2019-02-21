package main

//DbUser DataBase Username
var DbUser = "ravi"

//DbPassword Database password
var DbPassword = "hello123"

//DbName DataBase Name
var DbName = "rest_api_example"

func main() {
	a := App{}
	a.Initialize(DbUser, DbPassword, DbName)
	a.Run(":8080")
}
