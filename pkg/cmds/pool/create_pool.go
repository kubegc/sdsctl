package pool

import (
	"github.com/kube-stack/sdsctl/pkg/constant"
	"github.com/kube-stack/sdsctl/pkg/k8s"
	"github.com/kube-stack/sdsctl/pkg/utils"
	"github.com/kube-stack/sdsctl/pkg/virsh"
	"github.com/urfave/cli/v2"
	"os"
	"path/filepath"
	"strconv"
)

func NewCreatePoolCommand() *cli.Command {
	return &cli.Command{
		Name:      "create-pool",
		Usage:     "create kvm pool for kubestack",
		UsageText: "sdsctl [global options] create-pool [options]",
		Action:    createPool,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "pool",
				Usage: "name of pool",
			},
			&cli.StringFlag{
				Name:  "type",
				Usage: "storage pool type ",
				Value: "dir",
			},
			&cli.StringFlag{
				Name:  "url",
				Usage: "url of pool",
			},
			&cli.StringFlag{
				Name:  "content",
				Usage: "storage pool type",
				Value: "vmd",
			},
			&cli.StringFlag{
				Name:  "auto-start",
				Usage: "if auto-start pool",
				Value: "true",
			},
			&cli.StringFlag{
				Name:  "source-host",
				Usage: "remote storage server ip",
			},
			&cli.StringFlag{
				Name:  "source-path",
				Usage: "mount path of remote storage server",
			},
		},
	}
}

func createPool(ctx *cli.Context) error {
	logger := utils.GetLogger()
	autoStart, err := strconv.ParseBool(ctx.String("auto-start"))
	if err != nil {
		logger.Errorf("strconv.ParseBool auto-start err:%+v", err)
		return err
	}
	sourceHost, sourcePath := ctx.String("source-host"), ctx.String("source-path")
	pool, err := virsh.CreatePool(ctx.String("pool"), ctx.String("type"), ctx.String("url"), sourceHost, sourcePath)
	if err != nil {
		logger.Errorf("CreatePool err:%+v", err)
		virsh.DeletePool(ctx.String("pool"))
		return err
	}
	if err := virsh.AutoStartPool(ctx.String("pool"), autoStart); err != nil {
		logger.Errorf("AutoStartPool err:%+v", err)
		return err
	}
	// write content file
	contentPath := filepath.Join(ctx.String("url"), "content")
	var content = []byte(ctx.String("content"))
	os.WriteFile(contentPath, content, 0666)
	// update vmp
	ksgvr := k8s.NewKsGvr(constant.VMPS_Kind)
	flags := utils.ParseFlagMap(ctx)
	delete(flags, "auto-start")
	delete(flags, "source-host")
	delete(flags, "source-path")
	info, err := pool.GetInfo()
	if err != nil {
		logger.Errorf("GetInfo err:%+v", err)
		return err
	}
	extra := map[string]interface{}{
		"state":      constant.CRD_Pool_Active,
		"autostart":  autoStart,
		"capacity":   virsh.UniformBytes(info.Capacity),
		"sourceHost": sourceHost,
		"sourcePath": sourcePath,
	}
	flags = utils.MergeFlags(flags, extra)
	if err := ksgvr.Update(ctx.Context, constant.DefaultNamespace, ctx.String("pool"), constant.CRD_Pool_Key, flags); err != nil {
		virsh.DeletePool(ctx.String("pool"))
		return err
	}
	return nil
}
