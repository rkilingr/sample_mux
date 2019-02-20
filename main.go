package main

func main() {
	a := App{}
	a.Initialize("ravi", "hello123", "rest_api_example")
	a.Run(":8080")
}
