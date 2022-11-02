package image

import (
	"errors"
	"fmt"
	"github.com/kube-stack/sdsctl/pkg/utils"
	"github.com/kube-stack/sdsctl/pkg/virsh"
	"github.com/urfave/cli/v2"
)

func NewCreateDiskImageSnapshotCommand() *cli.Command {
	return &cli.Command{
		Name:      "create-disk-image",
		Usage:     "create disk image for kubestack",
		UsageText: "sdsctl [global options] create-disk-image [options]",
		Action:    createDiskImage,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "type",
				Usage: "storage vol type",
				Value: "dir",
			},
			&cli.StringFlag{
				Name:  "pool",
				Usage: "storage pool name",
			},
			&cli.StringFlag{
				Name:  "name",
				Usage: "storage volume disk name",
			},
			&cli.StringFlag{
				Name:  "source",
				Usage: "source storage image file",
			},
			&cli.StringFlag{
				Name:  "full-copy",
				Usage: "if full_copy, new disk will be created by snapshot",
			},
		},
	}
}

func createDiskImage(ctx *cli.Context) error {
	//logger := utils.GetLogger()
	pool := ctx.String("pool")
	active, err := virsh.IsPoolActive(pool)
	if err != nil {
		return err
	} else if !active {
		return fmt.Errorf("pool %+v is inactive", pool)
	}
	exist := virsh.IsDiskSnapshotExist(pool, ctx.String("source"), ctx.String("name"))
	if exist {
		return errors.New(fmt.Sprintf("the volume %+v is already exist", ctx.String("source")))
	}

	diskDir, _ := virsh.ParseDiskDir(pool, ctx.String("source"))
	config, err := virsh.ParseConfig(diskDir)
	if err != nil {
		return err
	}
	if !utils.Exists(config["current"]) {
		return errors.New(fmt.Sprintf("current disk %s not exists", config["current"]))
	}

	return nil
}
