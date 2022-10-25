package pool

import (
	"fmt"
	"github.com/kube-stack/sdsctl/pkg/constant"
	"github.com/kube-stack/sdsctl/pkg/k8s"
	"github.com/kube-stack/sdsctl/pkg/virsh"
	"github.com/urfave/cli/v2"
	"strconv"
)

func NewAutoStartPoolCommand() *cli.Command {
	return &cli.Command{
		Name:      "auto-start-pool",
		Usage:     "auto-start kvm pool for kubestack",
		UsageText: "sdsctl [global options] auto-start-pool [options]",
		Action:    autostartPool,
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
			&cli.StringFlag{
				Name:  "auto-start",
				Usage: "if auto-start pool",
				Value: "true",
			},
		},
	}
}

func autostartPool(ctx *cli.Context) error {
	autoStart, err := strconv.ParseBool(ctx.String("auto-start"))
	if err != nil {
		return err
	}
	if err := virsh.AutoStartPool(ctx.String("pool"), autoStart); err != nil {
		return err
	}
	// update vmp
	ksgvr := k8s.NewKsGvr(constant.VMPS_Kind)
	updateKey := fmt.Sprintf("%s.autostart", constant.CRD_Pool_Key)
	if err := ksgvr.Update(ctx.Context, constant.DefaultNamespace, ctx.String("pool"), updateKey, autoStart); err != nil {
		return err
	}
	return nil
}
