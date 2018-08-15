package main

import (
	app "github.com/dbaas-rest-api/app"
)

func main() {
	a := app.App{}
	a.Initialize("root", "", "", "mongodb://127.0.0.1:27017", 0)
	a.Run(":8081")
}
