package snapshot

import (
	"errors"
	"fmt"
	"github.com/kube-stack/sdsctl/pkg/constant"
	"github.com/kube-stack/sdsctl/pkg/k8s"
	"github.com/kube-stack/sdsctl/pkg/utils"
	"github.com/kube-stack/sdsctl/pkg/virsh"
	"github.com/urfave/cli/v2"
	"os"
	"path/filepath"
)

func NewDeleteExternalSnapshotCommand() *cli.Command {
	return &cli.Command{
		Name:      "delete-external-snapshot",
		Usage:     "delete kvm snapshot for kubestack",
		UsageText: "sdsctl [global options] delete-external-snapshot [options]",
		Action:    deleteExternalSnapshot,
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

func deleteExternalSnapshot(ctx *cli.Context) error {
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
	targetSSDir := filepath.Join(diskDir, "snapshots")
	targetSSPath := filepath.Join(targetSSDir, ctx.String("name"))
	backFile, _ := virsh.GetBackFile(targetSSPath)
	snapshotFiles := utils.GetFilesUnderDir(targetSSDir)
	files, err := virsh.GetBackChainFiles(snapshotFiles, targetSSPath)
	if err != nil {
		return err
	}
	// add snapshot
	files[ctx.String("name")] = true

	// 删除的是current的祖先
	vol := ctx.String("source")
	if _, ok := files[vol]; ok {
		vmActive, err := virsh.IsVMActive(domain)
		if err != nil {
			return err
		}
		// todo check?
		if domain != "" && vmActive {
			// live chain
			if err := virsh.LiveBlockForVMDisk(domain, config["current"], backFile); err != nil {
				return err
			}
		} else {
			if err := virsh.RebaseDiskSnapshot(backFile, config["current"], ""); err != nil {
				return err
			}
		}
	}

	delete(files, filepath.Base(config["current"]))
	// delete files
	for k, _ := range files {
		fullPath := filepath.Join(config["dir"], "snapshots", k)
		os.Remove(fullPath)
	}

	// delete vmdsn
	ksgvr := k8s.NewKsGvr(constant.VMDSNS_Kinds)
	for k, _ := range files {
		if err := ksgvr.Delete(ctx.Context, constant.DefaultNamespace, k); err != nil {
			return err
		}
	}
	return nil
}
