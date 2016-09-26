package get

import (
	"fmt"

	"github.com/urfave/cli"
)

type appsCmd struct{}

func (ac appsCmd) appsAction(ctx *cli.Context) error {
	fmt.Println("here")
	return nil
}
