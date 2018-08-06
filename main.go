// main.go

package main

func main() {
	a := App{}

	a.Initialize("root", "", "dbaas")

	a.Run(":8080")
}