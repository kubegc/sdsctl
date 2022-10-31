package virsh

import (
	"encoding/xml"
	"errors"
	"fmt"
	libvirtxml "github.com/libvirt/libvirt-go-xml"
)

func GetVMDiskSpec(domainName string) ([]libvirtxml.DomainDisk, error) {
	conn, err := GetConn()
	defer conn.Close()
	domain, err := conn.LookupDomainByName(domainName)
	// parse old format
	vxml, err := domain.GetXMLDesc(0)
	if err != nil {
		return nil, err
	}
	vmObj := &libvirtxml.Domain{}
	err = xml.Unmarshal([]byte(vxml), vmObj)
	if err != nil {
		return nil, err
	}
	disks := vmObj.Devices.Disks
	return disks, nil
}

func ParseVMDiskSpec(domainName string) (map[string]string, error) {
	disks, err := GetVMDiskSpec(domainName)
	if err != nil {
		return nil, err
	}
	res := make(map[string]string)
	for _, disk := range disks {
		if disk.Source != nil {
			res[disk.Source.File.File] = disk.Target.Dev
		}
	}
	return res, nil
}

func CheckVMDiskSpec(domainName, diskPath string) (map[string]string, error) {
	res, err := ParseVMDiskSpec(domainName)
	if err != nil {
		return res, err
	}
	for k, _ := range res {
		if k == diskPath {
			return res, nil
		}
	}
	return res, errors.New(fmt.Sprintf("domain %s has no disk %s", domainName, diskPath))
}
