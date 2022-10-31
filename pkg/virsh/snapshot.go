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
	if err := json.Unmarshal([]byte(output), res); err != nil {
		return "", err
	}
	return res["backing-filename"], nil
}
