package pool

import (
	"fmt"
	"github.com/WANNA959/sdsctl/pkg/constant"
	"github.com/WANNA959/sdsctl/pkg/k8s"
	"github.com/WANNA959/sdsctl/pkg/virsh"
	"github.com/dustin/go-humanize"
	"github.com/urfave/cli/v2"
)

func NewShowPoolCommand() *cli.Command {
	return &cli.Command{
		Name:      "show-pool",
		Usage:     "show kvm pool for kubestack",
		UsageText: "sdsctl [global options] show-pool [options]",
		Action:    showPool,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "type",
				Usage: "storage pool type",
				Value: "dir",
			},
			&cli.StringFlag{
				Name:  "pool",
				Usage: "storage pool type",
			},
		},
	}
}

func showPool(ctx *cli.Context) error {
	name := ctx.String("pool")
	pool, err := virsh.GetPoolInfo(name)
	if err != nil {
		return err
	}
	info, _ := pool.GetInfo()
	ksgvr := k8s.NewKsGvr(constant.VMPS_Kind)
	vmp, err := ksgvr.Get(ctx.Context, "default", name)
	if err != nil {
		return err
	}

	uuid, _ := pool.GetUUIDString()
	res := map[string]interface{}{
		"state":      virsh.GetPoolState(info.State),
		"uuid":       uuid,
		"free":       humanize.Bytes(info.Available),
		"capacity":   humanize.Bytes(info.Capacity),
		"allocation": humanize.Bytes(info.Allocation),
		"msg":        vmp.Spec.String(),
	}
	fmt.Printf("res: %+v", res)

	return nil
}
