package pool

import (
	"github.com/WANNA959/sdsctl/pkg/virsh"
	"github.com/urfave/cli/v2"
)

func NewAutoStartPoolCommand() *cli.Command {
	return &cli.Command{
		Name:      "auto-start-pool",
		Usage:     "auto-start kvm pool for kubestack",
		UsageText: "sdsctl [global options] auto-start-pool [options]",
		Action:    autostartPool,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "pool",
				Usage: "name of storage pool",
			},
			&cli.StringFlag{
				Name:  "type",
				Usage: "storage pool type ",
				Value: "dir",
			},
			&cli.BoolFlag{
				Name:  "auto-start",
				Usage: "if auto-start pool",
				Value: true,
			},
		},
	}
}

func autostartPool(ctx *cli.Context) error {
	return virsh.AutoStartPool(ctx.String("pool"), ctx.Bool("auto-start"))
}
