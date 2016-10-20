package main

import (
	"os"

	"github.com/feedhenry/rhm/commands"
	"github.com/feedhenry/rhm/commands/delete"
	"github.com/feedhenry/rhm/commands/get"
	"github.com/feedhenry/rhm/commands/use"
	"github.com/feedhenry/rhm/storage"

	"github.com/urfave/cli"
)

var outPutType string

//Note are using github.com/urfave/cli to do some common cli work
func main() {
	app := cli.NewApp()
	app.Name = "rhm"
	app.Usage = "a simple cli interface for Redhat Mobile Application Platform"
	app.Version = "0.0.1"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "o",
			Destination: &outPutType,
			Usage:       "-o=json",
			Value:       "plain",
		},
	}
	//create out data store for local file system
	store := storage.Store{}
	app.Commands = []cli.Command{
		//Login
		commands.NewLoginCmd(os.Stdin, os.Stdout, store),
		//Context
		commands.NewContextCmd(os.Stdin, os.Stdout, store),
		//Use
		use.NewUseCmd(os.Stdout, os.Stdin, store),
		//Get
		get.NewGetCmd(os.Stdout, os.Stdin, store),
		//delete
		delete.NewDeleteCmd(os.Stdin, os.Stdout, store),
	}
	app.Run(os.Args)
}
