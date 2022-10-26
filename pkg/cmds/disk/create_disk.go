package disk

import (
	"errors"
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/kube-stack/sdsctl/pkg/constant"
	"github.com/kube-stack/sdsctl/pkg/k8s"
	"github.com/kube-stack/sdsctl/pkg/utils"
	"github.com/kube-stack/sdsctl/pkg/virsh"
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
				Value: "localfs",
			},
			&cli.StringFlag{
				Name:  "vol",
				Usage: "storage vol name",
			},
			&cli.StringFlag{
				Name:  "capacity",
				Usage: "storage vol name",
			},
			&cli.StringFlag{
				Name:  "format",
				Usage: "storage vol format",
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
	exist := virsh.IsDiskExist(pool, ctx.String("vol"), ctx.String("type"))
	if exist {
		return errors.New(fmt.Sprintf("the volume %+v is already exist", ctx.String("vol")))
	}

	bytes, err := humanize.ParseBytes(ctx.String("capacity"))
	if err != nil {
		return err
	}

	if err = virsh.CreateDisk(pool, ctx.String("vol"), ctx.String("type"), ctx.String("capacity"), ctx.String("format")); err != nil {
		return err
	}
	// craete config.json
	diskPath, _ := virsh.ParseDiskDir(pool, ctx.String("vol"), ctx.String("type"))
	volPath, _ := virsh.ParseDiskPath(pool, ctx.String("vol"), ctx.String("type"), ctx.String("format"))
	cfg := map[string]string{
		"name":    ctx.String("vol"),
		"dir":     diskPath,
		"current": volPath,
		"pool":    pool,
	}
	if err := virsh.CreateConfig(diskPath, cfg); err != nil {
		return err
	}
	// update vmd
	ksgvr := k8s.NewKsGvr(constant.VMDS_Kind)
	flags := utils.ParseFlagMap(ctx)
	delete(flags, "capacity")
	extra := map[string]interface{}{
		"current":  volPath,
		"capacity": humanize.Bytes(bytes),
	}
	flags = utils.MergeFlags(flags, extra)
	if err = ksgvr.Update(ctx.Context, constant.DefaultNamespace, ctx.String("vol"), constant.CRD_Volume_Key, flags); err != nil {
		return err
	}
	return err
}
