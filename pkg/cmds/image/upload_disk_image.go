package image

import (
	"fmt"
	"github.com/kube-stack/sdsctl/pkg/k8s"
	"github.com/kube-stack/sdsctl/pkg/utils"
	"github.com/kube-stack/sdsctl/pkg/virsh"
	"github.com/urfave/cli/v2"
)

func NewUploadDiskImageCommand() *cli.Command {
	return &cli.Command{
		Name:      "upload-disk-image",
		Usage:     "upload disk image for kubestack",
		UsageText: "sdsctl [global options] upload-disk-image [options]",
		Action:    uploadDiskImage,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "source-path",
				Usage: "vmdi file path",
			},
			&cli.StringFlag{
				Name:  "pool",
				Usage: "vmdi storage pool name",
			},
			&cli.StringFlag{
				Name:  "vol",
				Usage: "storage volume disk image name",
			},
			&cli.StringFlag{
				Name:  "target-path",
				Usage: "nfs share path",
			},
		},
	}
}

func uploadDiskImage(ctx *cli.Context) error {
	logger := utils.GetLogger()
	pool := ctx.String("pool")
	uploadPath := ctx.String("source-path")
	targetPath := ctx.String("target-path")
	if uploadPath == "" {
		active, err := virsh.IsPoolActive(pool)
		if err != nil {
			return err
		} else if !active {
			return fmt.Errorf("pool %+v is inactive", pool)
		}
		if !virsh.CheckPoolType(pool, "vmdi") {
			return fmt.Errorf("pool type error")
		}
		uploadPath, err = virsh.ParseDiskPath(pool, ctx.String("vol"))
	}

	ip, err := k8s.GetNfsServiceIp()
	if err != nil {
		logger.Errorf("fail to get nfs service ip")
		return err
	}
	if !k8s.CheckNfsMount(ip, targetPath) {
		return fmt.Errorf("plz mount nfs path first")
	}

	return utils.CopyFile(uploadPath, targetPath)
}
