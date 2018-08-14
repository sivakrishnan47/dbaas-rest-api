// main.go

package main
import(
	app "./app"
)

func main() {
	a :=app.App{}

	a.Initialize("root", "", "")

	a.Run(":8080")
}