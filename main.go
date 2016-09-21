package main

import (
	"os"

	"github.com/feedhenry/rhm/commands"
	"github.com/feedhenry/rhm/commands/get"
	"github.com/feedhenry/rhm/storage"
	"github.com/urfave/cli"
)

//Note are using github.com/urfave/cli to do some common cli work
func main() {
	app := cli.NewApp()
	app.Name = "rhm"
	app.Version = "0.0.1"
	//create out data store for local file system
	store := storage.Store{}
	app.Commands = []cli.Command{
		commands.NewLoginCmd(os.Stdout, os.Stdin, store),
		get.NewGetCmd(os.Stdout, os.Stdin),
	}
	app.Run(os.Args)
}
