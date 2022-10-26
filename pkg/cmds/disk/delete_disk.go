package disk

import (
	"errors"
	"fmt"
	"github.com/kube-stack/sdsctl/pkg/constant"
	"github.com/kube-stack/sdsctl/pkg/k8s"
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
	exist := virsh.IsDiskExist(pool, ctx.String("vol"), ctx.String("type"))
	if !exist {
		return errors.New(fmt.Sprintf("the volume %+v is not exist", ctx.String("vol")))
	}

	if err = virsh.DeleteDisk(pool, ctx.String("vol"), ctx.String("type")); err != nil {
		return err
	}

	// delete vmd
	ksgvr := k8s.NewKsGvr(constant.VMDS_Kind)
	if ksgvr.Delete(ctx.Context, constant.DefaultNamespace, ctx.String("vol")); err != nil {
		return err
	}
	return nil
}
