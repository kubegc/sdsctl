package vm

import (
	"github.com/urfave/cli/v2"
)

func NewMigrateCommand() *cli.Command {
	return &cli.Command{
		Name:      "migrate",
		Usage:     "migrate vm for kubestack",
		UsageText: "sdsctl [global options] migrate [options]",
		Action:    migrate,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "domain",
				Usage: "kvm domain name",
			},
			&cli.StringFlag{
				Name:  "ip",
				Usage: "node ip",
			},
			&cli.StringFlag{
				Name:  "offline",
				Usage: "if support migrate offline",
			},
		},
	}
}

func migrate(ctx *cli.Context) error {
	// todo
	//domain := ctx.String("domain")

	return nil
}
