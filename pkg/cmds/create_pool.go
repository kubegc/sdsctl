package cmds

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"text/template"
)

var healthTemplate = template.Must(template.New("sdsctl create-pool").Parse(`
------------------------------------------------
network-controller:
    control grpc client health: {{.CtrlHealth}}
    bootstrap grpc client health: {{.BootHealth}}
------------------------------------------------
`))

func NewCreatePoolCommand() *cli.Command {
	return &cli.Command{
		Name:      "create-pool",
		Usage:     "create kvm pool for kubestack",
		UsageText: "sdsctl [global options] create-pool [options]",
		Action:    createPool,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "url",
				Usage: "url of pool",
			},
			&cli.StringFlag{
				Name:  "type",
				Usage: "storage pool type ",
				Value: "dir",
			},
			&cli.StringFlag{
				Name:  "content",
				Usage: "storage pool type",
				Value: "vmd",
			},
			&cli.BoolFlag{
				Name:  "auto-start",
				Usage: "if auto-start pool",
				Value: true,
			},
			&cli.StringFlag{
				Name:  "opt",
				Usage: "extra options",
			},
		},
	}
}

func createPool(ctx *cli.Context) error {
	ctype := ctx.String("type")
	fmt.Println(ctype)
	return nil
}
