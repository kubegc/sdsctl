package pool

import (
	"github.com/kube-stack/sdsctl/pkg/virsh"
	"github.com/urfave/cli/v2"
)

func NewCreatePoolCommand() *cli.Command {
	return &cli.Command{
		Name:      "create-pool",
		Usage:     "create kvm pool for kubestack",
		UsageText: "sdsctl [global options] create-pool [options]",
		Action:    createPool,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "name",
				Usage: "url of pool",
			},
			&cli.StringFlag{
				Name:  "type",
				Usage: "storage pool type ",
				Value: "dir",
			},
			&cli.StringFlag{
				Name:  "url",
				Usage: "url of pool",
			},
			&cli.StringFlag{
				Name:  "content",
				Usage: "storage pool type",
				Value: "vmd",
			},
			&cli.BoolFlag{
				Name:  "auto-start",
				Usage: "if auto-start pool",
				Value: true,
			},
			&cli.StringFlag{
				Name:  "opt",
				Usage: "extra options",
			},
		},
	}
}

func createPool(ctx *cli.Context) error {
	_, err := virsh.CreatePool(ctx.String("name"), ctx.String("type"), ctx.String("target"))
	if err != nil {
		return err
	}
	err = virsh.AutoStartPool(ctx.String("name"), ctx.Bool("auto-start"))
	return err
}
