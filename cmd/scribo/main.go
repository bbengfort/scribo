package main

import (
	"os"

	"github.com/bbengfort/scribo/scribo"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	addr := ":" + port

	app := scribo.CreateApp()
	app.Run(addr)
}
