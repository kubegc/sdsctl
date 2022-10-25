package disk

import (
	"errors"
	"fmt"
	"github.com/kube-stack/sdsctl/pkg/virsh"
	"github.com/urfave/cli/v2"
)

func NewCloneDiskCommand() *cli.Command {
	return &cli.Command{
		Name:      "clone-disk",
		Usage:     "clone kvm disk for kubestack",
		UsageText: "sdsctl [global options] clone-disk [options]",
		Action:    cloneDisk,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "type",
				Usage: "storage vol type",
				Value: "dir",
			},
			&cli.StringFlag{
				Name:  "pool",
				Usage: "storage pool name",
				Value: "dir",
			},
			&cli.StringFlag{
				Name:  "vol",
				Usage: "storage vol name",
			},
			&cli.StringFlag{
				Name:  "newvol",
				Usage: "new vol name",
			},
		},
	}
}

func cloneDisk(ctx *cli.Context) error {
	pool := ctx.String("pool")
	active, err := virsh.IsPoolActive(pool)
	if err != nil {
		return err
	} else if !active {
		return fmt.Errorf("pool %+v is inactive", pool)
	}
	exist := virsh.IsVolExist(pool, ctx.String("vol"), ctx.String("type"))
	if exist {
		return errors.New(fmt.Sprintf("the volume %+v is already exist", ctx.String("vol")))
	}

	return virsh.CloneVol(pool, ctx.String("vol"), ctx.String("newvol"), ctx.String("type"))
}
