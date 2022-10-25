package pool

import (
	"fmt"
	"github.com/kube-stack/sdsctl/pkg/constant"
	"github.com/kube-stack/sdsctl/pkg/k8s"
	"github.com/kube-stack/sdsctl/pkg/virsh"
	"github.com/urfave/cli/v2"
)

func NewStartPoolCommand() *cli.Command {
	return &cli.Command{
		Name:      "start-pool",
		Usage:     "start kvm pool for kubestack",
		UsageText: "sdsctl [global options] start-pool [options]",
		Action:    startPool,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "pool",
				Usage: "name of storage pool",
			},
			&cli.StringFlag{
				Name:  "type",
				Usage: "storage pool type ",
				Value: "dir",
			},
		},
	}
}

func startPool(ctx *cli.Context) error {
	if err := virsh.StartPool(ctx.String("pool")); err != nil {
		return err
	}
	// update vmp
	ksgvr := k8s.NewKsGvr(constant.VMPS_Kind)
	updateKey := fmt.Sprintf("%s.state", constant.CRD_Pool_Key)
	if err := ksgvr.Update(ctx.Context, constant.DefaultNamespace, ctx.String("pool"), updateKey, constant.CRD_Pool_Active); err != nil {
		return err
	}
	return nil
}
