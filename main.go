package main

func main() {
	a := App{}
	// You need to set your Username and Password here
	a.Initialize("DB_USERNAME", "DB_PASSWORD", "rest_api_example")

	a.Run(":8080")
}
