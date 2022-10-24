package virsh

import (
	"fmt"
	"testing"
)

func TestShowVol(t *testing.T) {
	vol, err := GetVol("pooltest2", "disktest", "dir")
	fmt.Printf("vol:%+v\n", vol)
	fmt.Printf("err:%+v\n", err)
}

func TestCreateVol(t *testing.T) {
	vol, err := CreateVol("pooltest2", "disktest.qcow2", "dir", "5G")
	fmt.Printf("vol:%+v\n", vol)
	fmt.Printf("err:%+v\n", err)
}

func TestDeleteVol(t *testing.T) {
	err := DeleteVol("pooltest2", "disktest.qcow2", "dir")
	fmt.Printf("err:%+v\n", err)
}
