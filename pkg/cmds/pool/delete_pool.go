package pool

import (
	"github.com/WANNA959/sdsctl/pkg/constant"
	"github.com/WANNA959/sdsctl/pkg/k8s"
	"github.com/WANNA959/sdsctl/pkg/virsh"
	"github.com/urfave/cli/v2"
)

func NewDeletePoolCommand() *cli.Command {
	return &cli.Command{
		Name:      "delete-pool",
		Usage:     "delete kvm pool for kubestack",
		UsageText: "sdsctl [global options] delete-pool [options]",
		Action:    deletePool,
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

func deletePool(ctx *cli.Context) error {
	ksgvr := k8s.NewKsGvr(constant.VMPS_Kind)
	err := ksgvr.Delete(ctx.Context, "default", ctx.String("pool"))
	if err != nil {
		return err
	}
	return virsh.DeletePool(ctx.String("pool"))
}
