package main

import (
	"fmt"
	"os"

	"github.com/bbengfort/scribo/scribo"
	"github.com/codegangsta/cli"
)

const version = "1.0"

func main() {

	// Instantiate the command line application
	app := cli.NewApp()
	app.Name = "scribo-register"
	app.Usage = "register a Node to Scribo and receive a Key for authentication."
	app.Version = version
	app.Author = "Benjamin Bengfort"
	app.Email = "benjamin@bengfort.com"
	app.Action = registerUser

	// Create the flags
	var addr string
	var dns string

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "addr",
			Value:       "",
			Usage:       "the IP address of the Node",
			Destination: &addr,
		},
		cli.StringFlag{
			Name:        "dns",
			Value:       "",
			Usage:       "the domain name of the Node",
			Destination: &dns,
		},
	}

	// Run the command line application
	app.Run(os.Args)
}

// The primary action of the scribo-register command
func registerUser(ctx *cli.Context) error {
	if ctx.NArg() == 1 {
		db := scribo.ConnectDB()
		name := ctx.Args()[0]

		// Get the Node out of the database
		node, err := scribo.GetNodeByName(db, name)

		// If the error is not nill, this is not simply does not exist
		if err != nil {
			cli.NewExitError(err.Error(), 2)
		}

		// Now we either have a new node or we are updating an old one.
		node.Name = name

		addr := ctx.String("addr")
		if addr != "" {
			node.Address = addr
		}

		dns := ctx.String("dns")
		if dns != "" {
			node.DNS = dns
		}

		// Reset the API key for the node.
		node.UpdateKey()

		// Now save the Node changes back to the database.
		created, err := node.Save(db)
		if err != nil {
			cli.NewExitError(err.Error(), 3)
		}

		// Print the node back to the console.
		var createStr string
		if created {
			createStr = "Created"
		} else {
			createStr = "Updated"
		}

		var addrStr string
		if node.Address != "" {
			addrStr = node.Address
		} else if node.DNS != "" {
			addrStr = node.DNS
		} else {
			addrStr = "Unknown Address"
		}

		fmt.Printf("%s Node %s (%s)\nKey: %s\n\n", createStr, node.Name, addrStr, node.Key)
		return nil

	} else if ctx.NArg() > 1 {
		return cli.NewExitError("Only one node can be registered at a time.", 1)
	}

	return cli.NewExitError("Supply the name of the node to register.", 1)
}
