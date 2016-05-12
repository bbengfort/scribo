package main

import (
	"os"
	"strconv"

	"github.com/bbengfort/scribo/scribo"
)

func main() {
	var port int
	var err error
	portEnv := os.Getenv("PORT")

	if portEnv == "" {
		port = 5356
	} else {
		port, err = strconv.Atoi(portEnv)
		if err != nil {
			panic(err)
		}
	}

	app := scribo.CreateApp()
	app.Run(port)
}
