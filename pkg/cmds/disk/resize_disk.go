package disk

import (
	"errors"
	"fmt"
	"github.com/kube-stack/sdsctl/pkg/virsh"
	"github.com/urfave/cli/v2"
)

func NewResizeDiskCommand() *cli.Command {
	return &cli.Command{
		Name:      "resize-disk",
		Usage:     "resize kvm disk for kubestack",
		UsageText: "sdsctl [global options] resize-disk [options]",
		Action:    resizeDisk,
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
				Name:  "capacity",
				Usage: "new storage vol capacity",
			},
		},
	}
}

func resizeDisk(ctx *cli.Context) error {
	pool := ctx.String("pool")
	active, err := virsh.IsPoolActive(pool)
	if err != nil {
		return err
	} else if !active {
		return fmt.Errorf("pool %+v is inactive", pool)
	}
	exist := virsh.IsDiskExist(pool, ctx.String("vol"))
	if !exist {
		return errors.New(fmt.Sprintf("the volume %+v is not exist", ctx.String("vol")))
	}

	return virsh.ResizeDisk(pool, ctx.String("vol"), ctx.String("capacity"))
}
