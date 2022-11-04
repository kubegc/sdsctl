package main

import (
	"fmt"
	"github.com/kube-stack/sdsctl/pkg/cmds"
	"github.com/kube-stack/sdsctl/pkg/cmds/disk"
	"github.com/kube-stack/sdsctl/pkg/cmds/externalSnapshot"
	"github.com/kube-stack/sdsctl/pkg/cmds/image"
	"github.com/kube-stack/sdsctl/pkg/cmds/internalSnapshot"
	"github.com/kube-stack/sdsctl/pkg/cmds/pool"
	"github.com/urfave/cli/v2"
	"os"
)

func main() {
	app := cmds.NewApp()
	app.Commands = []*cli.Command{
		// pool commands
		pool.NewShowPoolCommand(),
		pool.NewCreatePoolCommand(),
		pool.NewDeletePoolCommand(),
		pool.NewStopPoolCommand(),
		pool.NewAutoStartPoolCommand(),
		pool.NewStartPoolCommand(),

		// disk commands
		disk.NewShowDiskCommand(),
		disk.NewCreateDiskCommand(),
		disk.NewDeleteDiskCommand(),
		disk.NewCloneDiskCommand(),
		disk.NewResizeDiskCommand(),

		// disk external snapshot
		externalSnapshot.NewCreateExternalSnapshotCommand(),
		externalSnapshot.NewRevertExternalSnapshotCommand(),
		externalSnapshot.NewDeleteExternalSnapshotCommand(),

		// disk image
		image.NewCreateDiskFromImageCommand(),
		image.NewCreateDiskImageCommand(),
		image.NewDeleteDiskImageCommand(),

		// disk internal snapshot
		internalSnapshot.NewCreateInternalSnapshotCommand(),
		internalSnapshot.NewRevertInternalSnapshotCommand(),
		internalSnapshot.NewDeleteInternalSnapshotCommand(),
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Printf("error options: %s\n", err.Error())
		os.Exit(-1)
	}
}
