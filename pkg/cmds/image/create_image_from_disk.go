package image

import (
	"fmt"
	"github.com/kube-stack/sdsctl/pkg/virsh"
	"github.com/urfave/cli/v2"
)

func NewCreateImageFromDiskCommand() *cli.Command {
	return &cli.Command{
		Name:      "create-image-from-disk",
		Usage:     "create image from disk for kubestack",
		UsageText: "sdsctl [global options] create-image-from-disk [options]",
		Action:    createImageFromDisk,
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
				Name:  "vol",
				Usage: "source storage disk file path",
			},
			&cli.StringFlag{
				Name:  "target-pool",
				Usage: "vmdi storage pool name",
			},
			&cli.StringFlag{
				Name:  "name",
				Usage: "storage volume disk name",
			},
		},
	}
}

func createImageFromDisk(ctx *cli.Context) error {
	//logger := utils.GetLogger()
	pool := ctx.String("pool")
	targetPool := ctx.String("target-pool")
	active, err := virsh.IsPoolActive(pool)
	if err != nil {
		return err
	} else if !active {
		return fmt.Errorf("pool %+v is inactive", pool)
	}
	active2, err := virsh.IsPoolActive(targetPool)
	if err != nil {
		return err
	} else if !active2 {
		return fmt.Errorf("pool %+v is inactive", targetPool)
	}

	if !virsh.CheckPoolType(targetPool, "vmdi") {
		return fmt.Errorf("pool type error")
	}
	if !virsh.IsDiskExist(pool, ctx.String("vol")) {
		return fmt.Errorf("storage vol %s not exist", ctx.String("vol"))
	}
	sourceDiskdir, _ := virsh.ParseDiskDir(pool, ctx.String("vol"))
	config, _ := virsh.ParseConfig(sourceDiskdir)
	return createImage(ctx, config["current"], ctx.String("name"), targetPool)
}
