package get

import (
	"fmt"

	"github.com/urfave/cli"
)

type appsCmd struct{}

func (ac appsCmd) appsAction(ctx *cli.Context) error {
	switch ctx.GlobalString("o") {
	case "json":
		fmt.Println("{msg: \"here\"}")
	default:
		fmt.Println("here")
	}

	return nil
}
