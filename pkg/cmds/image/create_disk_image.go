package image

import (
	"fmt"
	"github.com/kube-stack/sdsctl/pkg/constant"
	"github.com/kube-stack/sdsctl/pkg/k8s"
	"github.com/kube-stack/sdsctl/pkg/utils"
	"github.com/kube-stack/sdsctl/pkg/virsh"
	"github.com/urfave/cli/v2"
	"os"
	"path/filepath"
	"strings"
)

func NewCreateDiskImageCommand() *cli.Command {
	return &cli.Command{
		Name:      "create-disk-image",
		Usage:     "create disk image for kubestack",
		UsageText: "sdsctl [global options] create-disk-image [options]",
		Action:    backcreateDiskImage,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "type",
				Usage: "storage vol type",
				Value: "dir",
			},
			&cli.StringFlag{
				Name:  "pool",
				Usage: "vmdi storage pool name",
			},
			&cli.StringFlag{
				Name:  "name",
				Usage: "storage volume disk image name",
			},
			&cli.StringFlag{
				Name:  "source",
				Usage: "source storage disk file path",
			},
		},
	}
}

func backcreateDiskImage(ctx *cli.Context) error {
	err := createDiskImage(ctx)
	ksgvr := k8s.NewKsGvr(constant.VMDIS_KINDS)
	if err != nil && !strings.Contains(err.Error(), "already exists") {
		ksgvr.UpdateWithStatus(ctx.Context, constant.DefaultNamespace, ctx.String("name"), constant.CRD_Volume_Key, nil, err.Error(), "400")
	}
	return err
}

func createImageBack(path string) {
	os.RemoveAll(path)
}

func createImage(ctx *cli.Context, sourceDiskPath, name, pool string) error {
	logger := utils.GetLogger()
	if !utils.Exists(sourceDiskPath) {
		return fmt.Errorf("disk file not exist")
	}

	targetImageDir, _ := virsh.ParseDiskDir(pool, name)
	if !utils.Exists(targetImageDir) {
		os.MkdirAll(targetImageDir, os.ModePerm)
	}
	targetImagePath := filepath.Join(targetImageDir, name)
	// cp source
	if err := utils.CopyFile(sourceDiskPath, targetImagePath); err != nil {
		return err
	}
	// rebase to self
	if err := virsh.RebaseDiskSnapshot("", targetImagePath, "qcow2"); err != nil {
		createImageBack(targetImageDir)
		logger.Errorf("RebaseDiskSnapshot err:%+v", err)
		return err
	}

	// write config
	cfg := map[string]string{
		"name":    name,
		"dir":     targetImageDir,
		"current": targetImagePath,
		"pool":    pool,
	}
	if err := virsh.CreateConfig(targetImageDir, cfg); err != nil {
		createImageBack(targetImageDir)
		logger.Errorf("CreateConfig err:%+v", err)
		return err
	}

	// create vmdi
	ksgvr := k8s.NewKsGvr(constant.VMDIS_KINDS)
	res := make(map[string]string)
	res["current"] = targetImagePath
	res["pool"] = pool
	res["format"] = "qcow2"
	res["type"] = ctx.String("type")
	if err := ksgvr.Create(ctx.Context, constant.DefaultNamespace, ctx.String("name"), constant.CRD_Volume_Key, res); err != nil {
		createImageBack(targetImageDir)
		logger.Errorf("ksgvr.Create err:%+v", err)
		return err
	}
	return nil
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
	if !virsh.CheckPoolType(pool, "vmdi") {
		return fmt.Errorf("pool type error")
	}
	return createImage(ctx, ctx.String("source"), ctx.String("name"), pool)
}
