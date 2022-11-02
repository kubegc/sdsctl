package disk

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

func NewCloneDiskCommand() *cli.Command {
	return &cli.Command{
		Name:      "clone-disk",
		Usage:     "clone kvm disk for kubestack",
		UsageText: "sdsctl [global options] clone-disk [options]",
		Action:    cloneDisk,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "type",
				Usage: "storage vol type",
				Value: "dir",
			},
			&cli.StringFlag{
				Name:  "pool",
				Usage: "storage pool name for newvol",
				Value: "dir",
			},
			&cli.StringFlag{
				Name:  "vol",
				Usage: "storage vol name",
			},
			&cli.StringFlag{
				Name:  "newvol",
				Usage: "new vol name",
			},
			&cli.StringFlag{
				Name:  "format",
				Usage: "new vol format",
			},
		},
	}
}

func cloneDisk(ctx *cli.Context) error {
	pool := ctx.String("pool")
	// new pool info & check
	poolGvr := k8s.NewKsGvr(constant.VMPS_Kind)
	vmp, err := poolGvr.Get(ctx.Context, constant.DefaultNamespace, ctx.String("pool"))
	if err != nil {
		return err
	}
	poolInfo, _ := k8s.GetCRDSpec(vmp.Spec.Raw, constant.CRD_Pool_Key)
	if poolInfo["state"] != constant.CRD_Pool_Active {
		return fmt.Errorf("pool %+v is inactive", pool)
	}

	// old disk info & check
	diskGvr := k8s.NewKsGvr(constant.VMDS_Kind)
	vmd, err := diskGvr.Get(ctx.Context, constant.DefaultNamespace, ctx.String("vol"))
	if err != nil {
		return err
	}
	sourceVolInfo, _ := k8s.GetCRDSpec(vmd.Spec.Raw, constant.CRD_Volume_Key)
	active, err := virsh.IsPoolActive(sourceVolInfo["pool"])
	if err != nil {
		return err
	} else if !active {
		return fmt.Errorf("pool %+v is inactive", sourceVolInfo["pool"])
	}
	exist := virsh.IsVolExist(sourceVolInfo["pool"], ctx.String("vol"))
	if exist {
		return errors.New(fmt.Sprintf("the volume %+v is not exist", ctx.String("vol")))
	}

	// path
	uuid := utils.GetUUID()
	middleDir := filepath.Join(poolInfo["url"], uuid)
	middlePath := filepath.Join(middleDir, ctx.String("newvol"))
	sourceDiskPath := filepath.Join(sourceVolInfo["current"])
	targetDiskDir := filepath.Join(poolInfo["url"], ctx.String("newvol"))
	targetDiskPath := filepath.Join(targetDiskDir, ctx.String("newvol"))

	// build middle disk with config file for mv or scp
	os.MkdirAll(middleDir, os.ModePerm)
	defer os.RemoveAll(middleDir)
	//if !utils.Exists(targetDiskDir) {
	//	os.MkdirAll(targetDiskDir, os.ModePerm)
	//}
	utils.CopyFile(sourceDiskPath, middlePath)
	file, _ := virsh.GetBackFile(middlePath)
	if file != "" {
		if err = virsh.RebaseDiskSnapshot("", middlePath, ctx.String("format")); err != nil {
			return err
		}
	}
	cfg := map[string]string{
		"name":    ctx.String("newvol"),
		"dir":     targetDiskDir,
		"current": targetDiskPath,
		"pool":    pool,
	}
	if err = virsh.CreateConfig(middleDir, cfg); err != nil {
		return err
	}

	// judge node name
	sourceNode := k8s.GetNodeName(vmd.Spec.Raw)
	targetNode := k8s.GetNodeName(vmp.Spec.Raw)
	if sourceNode == targetNode {
		// in same node
		if err := os.Rename(middleDir, targetDiskDir); err != nil {
			return err
		}
	} else {
		// in different node
		targetIP, _ := k8s.GetIPByNodeName(targetNode)
		if err = utils.CopyToRemoteFile(targetIP, middleDir, targetDiskDir); err != nil {
			return err
		}
		// todo remote register?
	}

	// create vmd
	ksgvr := k8s.NewKsGvr(constant.VMDS_Kind)
	sourceVmd, err := ksgvr.Get(ctx.Context, constant.DefaultNamespace, ctx.String("vol"))
	if err != nil {
		return err
	}
	ressourceVmdMap, _ := k8s.GetCRDSpec(sourceVmd.Spec.Raw, constant.CRD_Volume_Key)
	res := make(map[string]string)
	res["disk"] = ctx.String("vol")
	res["vol"] = ctx.String("vol")
	res["current"] = targetDiskPath
	res["pool"] = pool
	res["capacity"] = ressourceVmdMap["capacity"]
	res["format"] = ctx.String("format")
	res["type"] = ctx.String("type")

	if err = ksgvr.Create(ctx.Context, constant.DefaultNamespace, ctx.String("newvol"), constant.CRD_Volume_Key, res); err != nil {
		return err
	}
	return nil
}
