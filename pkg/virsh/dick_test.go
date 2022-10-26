package virsh

import (
	"fmt"
	"testing"
)

func TestCreateDisk(t *testing.T) {
	err := CreateDisk("pooltest2", "disktest2", "localfs", "10G", "qcow2")
	fmt.Printf("err: %+v", err)
}

func TestGetDisk(t *testing.T) {
	disk, err := GetDisk("pooltest2", "disktest2", "localfs", "qcow2")
	fmt.Printf("disk: %+v", disk)
	fmt.Printf("err: %+v", err)
}

func TestIsDiskExist(t *testing.T) {
	exist := IsDiskExist("pooltest2", "disktest2", "localfs")
	fmt.Printf("exist: %+v", exist)
}

func TestDeleteDisk(t *testing.T) {
	err := DeleteDisk("pooltest2", "disktest2", "localfs")
	fmt.Printf("err: %+v", err)
}
