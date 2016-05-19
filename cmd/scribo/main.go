// A command that runs the Scribo web API server
package main

import (
	"os"

	"github.com/bbengfort/scribo/scribo"
	"github.com/codegangsta/cli"
	"github.com/joho/godotenv"
)

func main() {
	// Load the .env file if it exists
	godotenv.Load()

	// Instantiate the command line application
	app := cli.NewApp()
	app.Name = "scribo"
	app.Usage = "runs the scribo web application and API"
	app.Version = scribo.Version
	app.Author = "Benjamin Bengfort"
	app.Email = "benjamin@bengfort.com"
	app.Action = runScriboApp

	// Create the flags
	var port int

	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:        "port",
			Value:       5356,
			Usage:       "the PORT to run the HTTP server on",
			EnvVar:      "PORT",
			Destination: &port,
		},
	}

	// Run the command line application
	app.Run(os.Args)
}

func runScriboApp(ctx *cli.Context) {
	server := scribo.CreateApp()
	server.Run(ctx.Int("port"))
}
