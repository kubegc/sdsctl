package internalSnapshot

import (
	"fmt"
	"github.com/kube-stack/sdsctl/pkg/virsh"
	"github.com/urfave/cli/v2"
)

func NewRevertInternalSnapshotCommand() *cli.Command {
	return &cli.Command{
		Name:      "revert-internal-snapshot",
		Usage:     "revert internal snapshot for kubestack",
		UsageText: "sdsctl [global options] revert-internal-snapshot [options]",
		Action:    revertInternalSnapshot,
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

func revertInternalSnapshot(ctx *cli.Context) error {
	//logger := utils.GetLogger()
	domainName := ctx.String("domain")
	if domainName == "" {
		return fmt.Errorf("domain can't be empty")
	}
	if !virsh.IsVMExist(domainName) {
		return fmt.Errorf("no domain named %s", domainName)
	}
	if err := checkDomainDisk(ctx, domainName); err != nil {
		return err
	}
	if err := virsh.RevertInternalSnapshot(domainName, ctx.String("snapshot")); err != nil {
		return err
	}
	return updateVMDSnapshot(ctx, domainName)
}
