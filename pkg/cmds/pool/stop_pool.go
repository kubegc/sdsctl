package pool

import (
	"github.com/kube-stack/sdsctl/pkg/virsh"
	"github.com/urfave/cli/v2"
)

func NewStopPoolCommand() *cli.Command {
	return &cli.Command{
		Name:      "stop-pool",
		Usage:     "stop kvm pool for kubestack",
		UsageText: "sdsctl [global options] stop-pool [options]",
		Action:    stopPool,
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

func stopPool(ctx *cli.Context) error {
	return virsh.StopPool(ctx.String("pool"))
}
