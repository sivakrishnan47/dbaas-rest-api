package main

import (
	"log"

	app "github.com/dbaas-rest-api/app"
)

func main() {
	log.Println("starting app")
	a := app.App{}
	log.Println("starting app after app")
	a.Initialize("root", "", "", "mongodb://172.17.0.3:27017", 0)
	a.Run(":8080")
}
