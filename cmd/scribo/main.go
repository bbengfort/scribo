package main

import "github.com/bbengfort/scribo/scribo"

func main() {
	app := scribo.CreateApp()
	app.Run("localhost:8080")
}
