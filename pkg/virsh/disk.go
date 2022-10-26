package virsh

import (
	"encoding/json"
	"fmt"
	"github.com/kube-stack/sdsctl/pkg/utils"
	"os"
	"path/filepath"
)

func ParseDiskDir(poolName, volName, vtype string) (string, error) {
	path, err := GetPoolTargetPath(poolName)
	if err != nil {
		return "", err
	}
	diskPath := filepath.Join(path, volName)
	return diskPath, err
}

func ParseDiskPath(poolName, volName, vtype, format string) (string, error) {
	diskPath, err := ParseDiskDir(poolName, volName, vtype)
	if err != nil {
		return "", err
	}
	volFile := fmt.Sprintf("%s.%s", volName, format)
	volPath := filepath.Join(diskPath, volFile)
	return volPath, nil
}

func GetDisk(poolName, volName, vtype, format string) (*Image, error) {
	volPath, err := ParseDiskPath(poolName, volName, vtype, format)
	if err != nil {
		return nil, err
	}
	image, err := OpenImage(volPath)
	return &image, err
}

func IsDiskExist(poolName, volName, vtype string) bool {
	diskPath, err := ParseDiskDir(poolName, volName, vtype)
	if err != nil {
		return false
	}
	return utils.IsDir(diskPath)
}

func CreateConfig(diskPath string, info map[string]string) error {
	// write content file
	configPath := filepath.Join(diskPath, "config.json")
	content, err := json.Marshal(info)
	if err != nil {
		return err
	}
	return os.WriteFile(configPath, content, 0666)
}

func CreateDisk(poolName, volName, vtype, capacity, format string) error {
	diskPath, err := ParseDiskDir(poolName, volName, vtype)
	if err != nil {
		return err
	}
	if !utils.Exists(diskPath) {
		os.MkdirAll(diskPath, os.ModePerm)
	}
	// create image
	volPath, err := ParseDiskPath(poolName, volName, vtype, format)
	if err != nil {
		return err
	}
	num, _ := parseCapacity(capacity)
	image := NewImage(volPath, format, num)
	return image.Create()
}

//func CreateDiskBack(poolName, volName, vtype, capacity, format string) error {
//
//	path, err := GetPoolTargetPath(poolName)
//	if err != nil {
//		return err
//	}
//	diskPath := filepath.Join(path, volName)
//	volFile := fmt.Sprintf("%s.%s", volName, format)
//	volPath := filepath.Join(diskPath, volFile)
//	cmd1 := &utils.Command{
//		Cmd: "mkdir -p " + volPath,
//	}
//	cmd2 := &utils.Command{
//		Cmd: "qemu-img create ",
//		Params: map[string]string{
//			"-f": "qcow2",
//			"":   fmt.Sprintf("%s %s", volPath, capacity),
//		},
//	}
//	cmds := utils.CommandList{
//		Comms: []*utils.Command{cmd1, cmd2},
//	}
//	return cmds.Execute()
//}

func DeleteDisk(poolName, volName, vtype string) error {
	path, err := GetPoolTargetPath(poolName)
	if err != nil {
		return err
	}
	diskPath := filepath.Join(path, volName)
	cmd := &utils.Command{
		Cmd: "rm -rf " + diskPath,
	}
	if _, err := cmd.Execute(); err != nil {
		return err
	}
	return nil
}
