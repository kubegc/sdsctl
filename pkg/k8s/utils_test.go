package k8s

import (
	"fmt"
	"testing"
)

func TestGetVMHostName(t *testing.T) {
	name := GetVMHostName()
	fmt.Println(name)
}

func TestGetIPByNodeName(t *testing.T) {
	ip, err := GetIPByNodeName("vm.wanna")
	fmt.Printf("err:%+v\n", err)
	fmt.Printf("name:%+v\n", ip)
}
