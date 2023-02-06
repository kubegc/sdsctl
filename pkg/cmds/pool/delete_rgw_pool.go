package pool

import (
	"github.com/kube-stack/sdsctl/pkg/constant"
	"github.com/kube-stack/sdsctl/pkg/k8s"
	"github.com/kube-stack/sdsctl/pkg/rook"
	"github.com/urfave/cli/v2"
)

func NewDeleteRgwPoolCommand() *cli.Command {
	return &cli.Command{
		Name:      "delete-rgw-pool",
		Usage:     "delete rgw image pool for kubestack",
		UsageText: "sdsctl [global options] delete-rgw-pool [options]",
		Action:    deleteRgwPool,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "name",
				Usage: "name of pool",
			},
		},
	}
}

func deleteRgwPool(ctx *cli.Context) error {
	name := ctx.String("name")
	// delete storageclass & obc
	if err := rook.DeleteOBC(name); err != nil {
		return err
	}
	if err := rook.DeleteBucketStorageClass(); err != nil {
		return err
	}

	// delete vmp
	ksgvr2 := k8s.NewKsGvr(constant.VMPS_Kind)
	return ksgvr2.Delete(ctx.Context, constant.DefaultNamespace, name)
}
