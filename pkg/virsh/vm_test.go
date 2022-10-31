package virsh

import (
	"fmt"
	"testing"
)

func TestGetVMDiskSpec(t *testing.T) {
	disks, err := GetVMDiskSpec("test")
	fmt.Printf("err:%+v\n", err)
	for _, disk := range disks {
		fmt.Printf("disk:%+v\n", disk)
	}
}
