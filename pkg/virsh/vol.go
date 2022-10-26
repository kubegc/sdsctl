package virsh

import (
	"encoding/xml"
	"fmt"
	"github.com/dustin/go-humanize"
	libvirtxml "github.com/libvirt/libvirt-go-xml"
	"libvirt.org/go/libvirt"
	"path/filepath"
	"strings"
)

func GetVol(poolName, volName, vtype string) (*libvirt.StorageVol, error) {
	conn, err := GetConn()
	defer conn.Close()
	if err != nil {
		return nil, err
	}
	pool, err := conn.LookupStoragePoolByName(poolName)
	if err != nil {
		return nil, err
	}
	vol, err := pool.LookupStorageVolByName(volName)
	if err != nil {
		return nil, err
	}
	return vol, nil
}

func IsVolExist(poolName, volName, vtype string) bool {
	_, err := GetVol(poolName, volName, vtype)
	if err != nil {
		return false
	}
	return true
}

var unitTransMap = map[string]string{
	"kb": "kib",
	"mb": "mib",
	"gb": "gib",
	"tb": "tib",
	"pb": "pib",
	"eb": "eib",
	"k":  "ki",
	"m":  "mi",
	"g":  "gi",
	"t":  "ti",
	"p":  "pi",
	"e":  "ei",
	"":   "b",
	"b":  "b",
}

func parseCapacity(raw string) (uint64, string) {
	var unit string
	var i int
	for i = len(raw) - 1; i >= 0; i-- {
		if raw[i] >= '0' && raw[i] <= '9' {
			unit = strings.ToLower(raw[i+1:])
			break
		}
	}
	if !strings.Contains(unit, "i") {
		unit = unitTransMap[unit]
	}
	newUnit := raw[0:i+1] + unit
	num, _ := humanize.ParseBytes(newUnit)
	return num, "byte"
}

func CreateVol(poolName, volName, vtype, capacity, format string) (*libvirt.StorageVol, error) {
	conn, err := GetConn()
	defer conn.Close()
	if err != nil {
		return nil, err
	}
	pool, err := conn.LookupStoragePoolByName(poolName)
	if err != nil {
		return nil, err
	}
	path, err := GetPoolTargetPath(poolName)
	if err != nil {
		return nil, err
	}
	diskPath := filepath.Join(path, volName)
	//if !utils.Exists(diskPath) {
	//	os.MkdirAll(diskPath, os.ModePerm)
	//}
	fmt.Println(diskPath)
	//pool, err = CreatePool(volName, vtype, diskPath)
	//if err != nil {
	//	return nil, err
	//}
	volFile := fmt.Sprintf("%s.%s", volName, format)
	volPath := filepath.Join(diskPath, volFile)
	fmt.Println(volPath)
	num, unit := parseCapacity(capacity)
	var volDesc = &libvirtxml.StorageVolume{
		Type: vtype,
		Name: volName,
		Capacity: &libvirtxml.StorageVolumeSize{
			Unit:  unit,
			Value: num,
		},
		Target: &libvirtxml.StorageVolumeTarget{
			Path: volPath,
			Format: &libvirtxml.StorageVolumeTargetFormat{
				Type: format,
			},
		},
	}
	volXML, err := volDesc.Marshal()
	if err != nil {
		return nil, err
	}
	return pool.StorageVolCreateXML(volXML, 0)
}

func DeleteVol(poolName, volName, vtype string) error {
	vol, err := GetVol(poolName, volName, vtype)
	if err != nil {
		return err
	}
	vol.Delete(0)
	return vol.Free()
}

func ResizeVol(poolName, volName, vtype, capacity string) error {
	vol, err := GetVol(poolName, volName, vtype)
	if err != nil {
		return err
	}
	bytes, _ := parseCapacity(capacity)
	return vol.Resize(bytes, libvirt.STORAGE_VOL_RESIZE_SHRINK)
}

func CloneVol(poolName, volName, newVolName, vtype string) error {
	pool, err := GetPoolInfo(poolName)
	if err != nil {
		return err
	}
	vol, err := GetVol(poolName, volName, vtype)
	if err != nil {
		return err
	}

	// parse old format
	vxml, err := vol.GetXMLDesc(0)
	if err != nil {
		return err
	}
	volObj := &libvirtxml.StorageVolume{}
	err = xml.Unmarshal([]byte(vxml), volObj)
	if err != nil {
		return err
	}
	format := volObj.Target.Format.Type

	// define vol xml, with newname & format
	xml := &libvirtxml.StorageVolume{
		Name: newVolName,
		Target: &libvirtxml.StorageVolumeTarget{
			Format: &libvirtxml.StorageVolumeTargetFormat{
				Type: format,
			},
		},
	}
	xmlStr, err := xml.Marshal()
	if err != nil {
		return err
	}
	// clone with old vol
	_, err = pool.StorageVolCreateXMLFrom(xmlStr, vol, 0)
	return err
}
