package virsh

import (
	"encoding/json"
	"fmt"
	"github.com/kube-stack/sdsctl/pkg/utils"
)

func CreateExternalSnapshot(domain, snapshot, diskPath, targetSSPath, noNeedSnapshotDisk, format string) error {
	cmd := utils.Command{
		Cmd: "virsh snapshot-create-as",
		Params: map[string]string{
			"--domain":         domain,
			"--name":           snapshot,
			"--atomic":         "",
			"--disk-only":      "",
			"--no-metadata":    "",
			"--diskspec":       fmt.Sprintf("%s,snapshot=external,file=%s,driver=%s", diskPath, targetSSPath, format),
			noNeedSnapshotDisk: "",
		},
	}
	if _, err := cmd.Execute(); err != nil {
		return err
	}
	return nil
}

func GetBackFile(path string) (string, error) {
	cmd := &utils.Command{
		Cmd: fmt.Sprintf("qemu-img info -U --output json %s", path),
	}
	output, err := cmd.Execute()
	if err != nil {
		return "", err
	}
	res := make(map[string]string)
	if err := json.Unmarshal([]byte(output), &res); err != nil {
		return "", err
	}
	if _, ok := res["backing-filename"]; !ok {
		return "", nil
	}
	return res["backing-filename"], nil
}

func GetBackChainFiles(snapshotFiles []string, backFile string) (map[string]bool, error) {
	res := make(map[string]bool)
	for _, snapshotFile := range snapshotFiles {
		files, err := GetOneBackChainFiles(snapshotFile)
		if err != nil {
			return res, err
		}
		if _, ok := files[backFile]; ok {
			res[snapshotFile] = true
		}
	}
	return res, nil
}

func GetOneBackChainFiles(path string) (map[string]bool, error) {
	res := make(map[string]bool)
	cmd := &utils.Command{
		Cmd: fmt.Sprintf("qemu-img info -U --backing-chain --output json %s", path),
	}
	output, err := cmd.ExecuteWithPlain()
	if err != nil {
		return res, err
	}
	infos := make([]map[string]interface{}, 0)
	if err := json.Unmarshal(output, &infos); err != nil {
		return res, err
	}
	for _, info := range infos {
		file, ok := info["backing-filename"]
		if ok {
			res[file.(string)] = true
		}
	}
	return res, nil
}

func LiveBlockForVMDisk(domainName, path, base string) error {
	parseBase := "''"
	if base != "" {
		parseBase = base
	}
	cmd := utils.Command{
		Cmd: "virsh blockpull",
		Params: map[string]string{
			"--domain": domainName,
			"--path":   path,
			"--base":   parseBase,
			"--wait":   "",
		},
	}
	_, err := cmd.Execute()
	return err
}

func RebaseDiskSnapshot(base, path string) error {
	parseBase := "''"
	if base != "" {
		parseBase = base
	}
	cmd := utils.Command{
		Cmd: "qemu-img rebase -b %s %s",
		Params: map[string]string{
			"-b": fmt.Sprintf("%s %s", parseBase, path),
		},
	}
	_, err := cmd.Execute()
	return err
}
