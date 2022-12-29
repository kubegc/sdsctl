package pool

import (
	"fmt"
	"github.com/kube-stack/sdsctl/pkg/constant"
	"github.com/kube-stack/sdsctl/pkg/k8s"
	"github.com/kube-stack/sdsctl/pkg/rook"
	"github.com/urfave/cli/v2"
)

func NewCreateNFSPoolCommand() *cli.Command {
	return &cli.Command{
		Name:      "create-nfs-pool",
		Usage:     "create nfs image pool for kubestack",
		UsageText: "sdsctl [global options] create-nfs-pool [options]",
		Action:    createNFSPool,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "name",
				Usage: "name of pool",
			},
			&cli.StringFlag{
				Name:  "local-path",
				Usage: "local mount path",
			},
		},
	}
}

func createNFSPool(ctx *cli.Context) error {
	name := ctx.String("name")
	//if err := rook.CreateNfsPool(name); err != nil {
	//	return err
	//}
	if err := rook.WaitNFSPoolReady(constant.DefaultNFSClusterName); err != nil {
		return err
	}
	nfsPath := name
	if err := rook.ExportNFSPath(constant.DefaultNFSClusterName, nfsPath); err != nil {
		return err
	}
	ip := rook.GetNfsServiceIp(constant.DefaultNFSClusterName)
	if len(ip) == 0 {
		return fmt.Errorf("fail to get nfs server ip")
	}
	if err := rook.MountNfs(ip, nfsPath, ctx.String("local-path")); err != nil {
		return err
	}
	//ksgvr := k8s.NewExternalGvr(constant.DefaultRookGroup, constant.DefaultRookVersion, constant.CephNFSPoolS_Kinds)
	ksgvr := k8s.NewKsGvr(constant.VMPS_Kind)
	flags := map[string]string{
		"pool":        name,
		"content":     "vmdi",
		"type":        constant.PoolCephNFSType,
		"server-path": fmt.Sprintf("%s:/%s", ip, nfsPath),
		"local-path":  ctx.String("local-path"),
	}
	if err := ksgvr.Update(ctx.Context, constant.RookNamespace, name, constant.CRD_Pool_Key, flags); err != nil {
		return err
	}
	return nil
}
