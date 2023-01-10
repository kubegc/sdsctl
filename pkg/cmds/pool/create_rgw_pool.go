package pool

import (
	"fmt"
	"github.com/kube-stack/sdsctl/pkg/constant"
	"github.com/kube-stack/sdsctl/pkg/k8s"
	"github.com/kube-stack/sdsctl/pkg/rook"
	"github.com/urfave/cli/v2"
)

func NewCreateRgwPoolCommand() *cli.Command {
	return &cli.Command{
		Name:      "create-rgw-pool",
		Usage:     "create rgw image pool for kubestack",
		UsageText: "sdsctl [global options] create-rgw-pool [options]",
		Action:    createRgwPool,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "name",
				Usage: "name of pool",
			},
		},
	}
}

func createRgwPool(ctx *cli.Context) error {
	name := ctx.String("name")
	// create storageclass & obc
	if err := rook.CreateBucketStorageClass(); err != nil {
		return err
	}
	if err := rook.CreateOBC(constant.DefaultCephRwgName); err != nil {
		return err
	}

	// update vmp crd
	info, err := rook.GetBucketInfo(constant.DefaultCephRwgName)
	if err != nil {
		return err
	}
	secret, err := rook.GetBucketSecret(constant.DefaultCephRwgName)
	if err != nil {
		return err
	}
	ksgvr := k8s.NewKsGvr(constant.VMPS_Kind)
	flags := map[string]string{
		"pool":           name,
		"content":        "vmdi",
		"type":           constant.PoolCephRgwType,
		"host":           info["host"],
		"server-address": fmt.Sprintf("%s:%s", info["ip"], info["port"]),
		"bucket":         info["name"],
		"access-id":      secret["access-id"],
		"access-key":     secret["access-key"],
	}
	if err := ksgvr.Update(ctx.Context, constant.RookNamespace, name, constant.CRD_Pool_Key, flags); err != nil {
		return err
	}
	return nil
}
