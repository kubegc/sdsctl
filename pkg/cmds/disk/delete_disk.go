package disk

import (
	"errors"
	"fmt"
	"github.com/kube-stack/sdsctl/pkg/virsh"
	"github.com/urfave/cli/v2"
)

func NewDeleteDiskCommand() *cli.Command {
	return &cli.Command{
		Name:      "delete-disk",
		Usage:     "delete kvm disk for kubestack",
		UsageText: "sdsctl [global options] delete-disk [options]",
		Action:    deleteDisk,
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
		},
	}
}

func deleteDisk(ctx *cli.Context) error {
	pool := ctx.String("pool")
	active, err := virsh.IsPoolActive(pool)
	if err != nil {
		return err
	} else if !active {
		return fmt.Errorf("pool %+v is inactive", pool)
	}
	exist := virsh.IsVolExist(pool, ctx.String("vol"), ctx.String("type"))
	if !exist {
		return errors.New(fmt.Sprintf("the volume %+v is not exist", ctx.String("vol")))
	}

	return virsh.DeleteVol(pool, ctx.String("vol"), ctx.String("type"))
}
