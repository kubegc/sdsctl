package disk

import (
	"errors"
	"fmt"
	"github.com/kube-stack/sdsctl/pkg/constant"
	"github.com/kube-stack/sdsctl/pkg/k8s"
	"github.com/kube-stack/sdsctl/pkg/utils"
	"github.com/kube-stack/sdsctl/pkg/virsh"
	"github.com/urfave/cli/v2"
	"path/filepath"
	"strconv"
)

func NewCreateDiskFromImageCommand() *cli.Command {
	return &cli.Command{
		Name:      "create-disk-from-image",
		Usage:     "create kvm disk from image for kubestack",
		UsageText: "sdsctl [global options] create-disk-from-image [options]",
		Action:    createDiskFromImage,
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
				Name:  "target",
				Usage: "new storage vol name",
			},
			&cli.StringFlag{
				Name:  "source",
				Usage: "source storage vol path",
			},
			&cli.StringFlag{
				Name:  "full-copy",
				Usage: "if full-copy, new disk will be created by snapshot",
			},
		},
	}
}

func createDiskFromImage(ctx *cli.Context) error {
	pool := ctx.String("pool")
	parseBool, err := strconv.ParseBool(ctx.String("full-copy"))
	if err != nil {
		return err
	}
	active, err := virsh.IsPoolActive(pool)
	if err != nil {
		return err
	} else if !active {
		return fmt.Errorf("pool %+v is inactive", pool)
	}
	exist := virsh.IsDiskExist(pool, ctx.String("source"), ctx.String("type"))
	if exist {
		return errors.New(fmt.Sprintf("the volume %+v is already exist", ctx.String("source")))
	}

	// source
	image, _ := virsh.OpenImage(ctx.String("source"))
	sourceFormat := image.Format
	// target
	targetDiskDir, _ := virsh.ParseDiskDir(pool, ctx.String("target"), ctx.String("type"))
	targetDiskPath := filepath.Join(targetDiskDir, ctx.String("target"))
	targetDiskConfig := filepath.Join(targetDiskDir, "config.json")
	if utils.Exists(targetDiskConfig) {
		return errors.New(fmt.Sprintf("target disk %s already exists", targetDiskConfig))
	}

	if parseBool {
		if err := virsh.CreateFullCopyDisk(ctx.String("source"), sourceFormat, targetDiskPath); err != nil {
			return err
		}
	} else {
		if err := virsh.CreateQcow2DiskWithBacking(ctx.String("source"), sourceFormat, targetDiskPath); err != nil {
			return err
		}
	}
	cfg := map[string]string{
		"name":    ctx.String("target"),
		"dir":     targetDiskDir,
		"current": targetDiskPath,
		"pool":    pool,
	}
	if err = virsh.CreateConfig(targetDiskDir, cfg); err != nil {
		return err
	}

	// update vmd
	ksgvr := k8s.NewKsGvr(constant.VMDS_Kind)
	flags := utils.ParseFlagMap(ctx)
	delete(flags, "target")
	extra := map[string]interface{}{
		"current": targetDiskPath,
	}
	flags = utils.MergeFlags(flags, extra)
	if err = ksgvr.Update(ctx.Context, constant.DefaultNamespace, ctx.String("target"), constant.CRD_Volume_Key, flags); err != nil {
		return err
	}
	return err
}
