package cmds

import (
	"fmt"
	"github.com/urfave/cli/v2"
)

func NewShowPoolCommand() *cli.Command {
	return &cli.Command{
		Name:      "show-pool",
		Usage:     "show kvm pool for kubestack",
		UsageText: "sdsctl [global options] show-pool [options]",
		Action:    showPool,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "type",
				Usage: "storage pool type",
				Value: "dir",
			},
			&cli.StringFlag{
				Name:  "pool",
				Usage: "storage pool type",
				Value: "vmd",
			},
		},
	}
}

func showPool(ctx *cli.Context) error {
	ctype := ctx.String("type")
	fmt.Println(ctype)
	return nil
}
