package disk

import (
	"errors"
	"fmt"
	"github.com/WANNA959/sdsctl/pkg/virsh"
	"github.com/urfave/cli/v2"
)

func NewCreateDiskCommand() *cli.Command {
	return &cli.Command{
		Name:      "create-disk",
		Usage:     "create kvm disk for kubestack",
		UsageText: "sdsctl [global options] create-disk [options]",
		Action:    createDisk,
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
				Usage: "storage vol name",
			},
		},
	}
}

func createDisk(ctx *cli.Context) error {
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

	_, err = virsh.CreateVol(pool, ctx.String("vol"), ctx.String("type"), ctx.String("capacity"), ctx.String("format"))
	return err
}
