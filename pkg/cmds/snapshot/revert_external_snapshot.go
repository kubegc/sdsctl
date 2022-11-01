package snapshot

import (
	"errors"
	"fmt"
	"github.com/kube-stack/sdsctl/pkg/constant"
	"github.com/kube-stack/sdsctl/pkg/k8s"
	"github.com/kube-stack/sdsctl/pkg/utils"
	"github.com/kube-stack/sdsctl/pkg/virsh"
	"github.com/urfave/cli/v2"
	"path/filepath"
)

func NewRevertExternalSnapshotCommand() *cli.Command {
	return &cli.Command{
		Name:      "revert-external-snapshot",
		Usage:     "revert kvm snapshot for kubestack",
		UsageText: "sdsctl [global options] revert-external-snapshot [options]",
		Action:    revertExternalSnapshot,
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
				Usage: "storage volume snapshot name",
			},
			&cli.StringFlag{
				Name:  "format",
				Usage: "storage vol format",
			},
			&cli.StringFlag{
				Name:  "source",
				Usage: "source storage disk file",
			},
			&cli.StringFlag{
				Name:  "domain",
				Usage: "domain name",
			},
		},
	}
}

// revert snapshot {name} 到上一个版本（back file）
func revertExternalSnapshot(ctx *cli.Context) error {
	domain := ctx.String("domain")
	pool := ctx.String("pool")
	active, err := virsh.IsPoolActive(pool)
	if err != nil {
		return err
	} else if !active {
		return fmt.Errorf("pool %+v is inactive", pool)
	}
	exist := virsh.IsDiskSnapshotExist(pool, ctx.String("source"), ctx.String("snapshot"))
	if !exist {
		return errors.New(fmt.Sprintf("the snapshot %+v is not exist", ctx.String("source")))
	}
	diskDir, _ := virsh.ParseDiskDir(pool, ctx.String("source"))
	config, err := virsh.ParseConfig(diskDir)
	if err != nil {
		return err
	}
	if !virsh.CheckDiskInUse(config["current"]) {
		return errors.New("current disk in use, plz check or set real domain field")
	}
	backFile, err := virsh.GetBackFile(config["current"])
	if err != nil {
		return err
	}
	newFile := utils.GetUUID()
	newFilePath := filepath.Join(utils.GetDir(backFile), newFile)
	if err := virsh.CreateDiskWithBacking(ctx.String("format"), backFile, ctx.String("format"), newFilePath); err != nil {
		return err
	}
	// change vm disk
	if domain != "" {
		if err := virsh.ChangeVMDisk(domain, config["current"], newFilePath); err != nil {
			return err
		}
	}

	// write config: current point to snapshot
	config["current"] = newFilePath
	virsh.CreateConfig(diskDir, config)

	// update vmd
	ksgvr := k8s.NewKsGvr(constant.VMDS_Kind)
	vmd, err := ksgvr.Get(ctx.Context, constant.DefaultNamespace, ctx.String("source"))
	if err != nil {
		return err
	}
	res, _ := k8s.GetCRDSpec(vmd.Spec.Raw, constant.CRD_Volume_Key)
	res["disk"] = ctx.String("source")
	res["current"] = newFilePath
	res["full_backing_filename"] = backFile
	// todo lifecycle?
	if err = ksgvr.Update(ctx.Context, constant.DefaultNamespace, ctx.String("source"), constant.CRD_Volume_Key, res); err != nil {
		return err
	}

	return nil
}
