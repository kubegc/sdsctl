package cmds

import (
	"fmt"
	"github.com/kube-stack/sdsctl/pkg/version"
	"github.com/urfave/cli/v2"
	"os"
	"path/filepath"
	"runtime"
)

type AccessConfig struct {
	ip            string
	port          string
	bootStrapPort string
	CAFile        string
	CertFile      string
	KeyFile       string
}

var GlobalConfig AccessConfig

var homeDir string = func() string {
	if home, err := os.UserHomeDir(); err != nil {
		return ""
	} else {
		return home
	}
}()

func NewApp() *cli.App {
	app := cli.NewApp()
	app.Name = "sdsctl"
	app.Usage = "sdsctl, a commond-line tool to control storage for kubestack"
	app.Version = version.Version
	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Printf("%s version %s\n", app.Name, app.Version)
		fmt.Printf("go version %s\n", runtime.Version())
	}
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:        "ip",
			Usage:       "leader host ip",
			Destination: &GlobalConfig.ip,
			Value:       "127.0.0.1",
		},
		&cli.StringFlag{
			Name:        "port",
			Usage:       "network grpc control port",
			Destination: &GlobalConfig.port,
			Value:       "6440",
		},
		&cli.StringFlag{
			Name:        "bootport",
			Usage:       "network grpc bootstrap control port",
			Destination: &GlobalConfig.bootStrapPort,
			Value:       "6439",
		},
		// todo defualt cert file path
		&cli.StringFlag{
			Name:        "cacert",
			Usage:       "ca cert filepath of network grpc server",
			Destination: &GlobalConfig.CAFile,
			Value:       filepath.Join(homeDir, ".litekube/nc/certs/grpc/ca.pem"),
		},
		&cli.StringFlag{
			Name:        "cert",
			Usage:       "client cert filepath of network grpc server",
			Destination: &GlobalConfig.CertFile,
			Value:       filepath.Join(homeDir, ".litekube/nc/certs/grpc/client.pem"),
		},
		&cli.StringFlag{
			Name:        "key",
			Usage:       "client key filepath of network grpc server",
			Destination: &GlobalConfig.KeyFile,
			Value:       filepath.Join(homeDir, ".litekube/nc/certs/grpc/client-key.pem"),
		},
	}

	return app
}
