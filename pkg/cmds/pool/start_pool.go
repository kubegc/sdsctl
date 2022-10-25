package pool

import (
	"github.com/kube-stack/sdsctl/pkg/virsh"
	"github.com/urfave/cli/v2"
)

func NewStartPoolCommand() *cli.Command {
	return &cli.Command{
		Name:      "start-pool",
		Usage:     "start kvm pool for kubestack",
		UsageText: "sdsctl [global options] start-pool [options]",
		Action:    startPool,
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
		},
	}
}

func startPool(ctx *cli.Context) error {
	return virsh.StartPool(ctx.String("pool"))
}
