package internalSnapshot

import (
	"fmt"
	"github.com/kube-stack/sdsctl/pkg/constant"
	"github.com/kube-stack/sdsctl/pkg/k8s"
	"github.com/kube-stack/sdsctl/pkg/utils"
	"github.com/kube-stack/sdsctl/pkg/virsh"
	"github.com/urfave/cli/v2"
	"strings"
)

func NewCreateInternalSnapshotCommand() *cli.Command {
	return &cli.Command{
		Name:      "create-internal-snapshot",
		Usage:     "create internal snapshot for kubestack",
		UsageText: "sdsctl [global options] create-internal-snapshot [options]",
		Action:    backcreateInternalSnapshot,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "type",
				Usage: "storage vol type",
				Value: "dir",
			},
			&cli.StringFlag{
				Name:  "domain",
				Usage: "domain name",
			},
			&cli.StringFlag{
				Name:  "name",
				Usage: "storage disk name",
			},
			&cli.StringFlag{
				Name:  "snapshot",
				Usage: "storage vol format",
			},
		},
	}
}

func backcreateInternalSnapshot(ctx *cli.Context) error {
	err := createInternalSnapshot(ctx)
	ksgvr := k8s.NewKsGvr(constant.VMDS_Kind)
	updateKey := fmt.Sprintf("%s.snapshots", constant.CRD_Volume_Key)
	if err != nil && !strings.Contains(err.Error(), "already exists") {
		ksgvr.UpdateWithStatus(ctx.Context, constant.DefaultNamespace, ctx.String("name"), updateKey, nil, err.Error(), "400")
	}
	return err
}

func checkDomainDisk(ctx *cli.Context, domainName string) error {
	ksgvr := k8s.NewKsGvr(constant.VMDS_Kind)
	vmd, err := ksgvr.Get(ctx.Context, constant.DefaultNamespace, ctx.String("name"))
	if err != nil {
		return err
	}
	spec, _ := k8s.GetCRDSpec(vmd.Spec.Raw, constant.CRD_Volume_Key)
	if _, err := virsh.CheckVMDiskSpec(domainName, spec["current"]); err != nil {
		return err
	}
	return nil
}

func updateVMDSnapshot(ctx *cli.Context, domainName string) error {
	// update vmd
	ksgvr := k8s.NewKsGvr(constant.VMDS_Kind)
	updateKey := fmt.Sprintf("%s.snapshots", constant.CRD_Volume_Key)
	snapshots, err := virsh.ListAllCurrentInternalSnapshots(domainName)
	if err != nil {
		return err
	}
	return ksgvr.Update(ctx.Context, constant.DefaultNamespace, ctx.String("name"), updateKey, snapshots)
}

func createInternalSnapshot(ctx *cli.Context) error {
	logger := utils.GetLogger()
	domainName := ctx.String("domain")
	if domainName == "" {
		return fmt.Errorf("domain can't be empty")
	}
	if !virsh.IsVMExist(domainName) {
		return fmt.Errorf("no domain named %s", domainName)
	}
	if err := checkDomainDisk(ctx, domainName); err != nil {
		logger.Errorf("checkDomainDisk err:%+v", err)
		return err
	}
	if err := virsh.CreateInternalSnapshot(domainName, ctx.String("snapshot")); err != nil {
		logger.Errorf("CreateInternalSnapshot err:%+v", err)
		return err
	}

	return updateVMDSnapshot(ctx, domainName)
}
