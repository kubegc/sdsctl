package utils

import (
	"fmt"
	"testing"
)

func TestGetDir(t *testing.T) {
	dir := GetDir("/var/lib/libvirt/pooltest2/disktest2/disktest2.qcow2")
	fmt.Println(dir)
}
