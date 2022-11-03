package externalSnapshot

import (
	"github.com/urfave/cli/v2"
)

func NewShowDiskSnapshotCommand() *cli.Command {
	return &cli.Command{
		Name:      "show-disk-snapshot",
		Usage:     "create kvm snapshot for kubestack",
		UsageText: "sdsctl [global options] show-disk-snapshot [options]",
		Action:    showDiskSnapshot,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "type",
				Usage: "storage vol type",
				Value: "dir",
			},
			&cli.StringFlag{
				Name:  "pool",
				Usage: "storage pool name",
			},
			&cli.StringFlag{
				Name:  "vol",
				Usage: "storage volume name",
			},
			&cli.StringFlag{
				Name:  "name",
				Usage: "storage volume snapshot name",
			},
		},
	}
}

func showDiskSnapshot(ctx *cli.Context) error {
	return nil
}
