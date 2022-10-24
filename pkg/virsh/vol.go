package virsh

import (
	libvirtxml "github.com/libvirt/libvirt-go-xml"
	"libvirt.org/go/libvirt"
	"path/filepath"
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

func parseCapacity(raw string) (uint64, string) {
	var num uint64 = 0
	var unit string
	for idx := range raw {
		if raw[idx] >= '0' && raw[idx] <= '9' {
			num = 10*num + uint64(raw[idx]-'0')
		} else {
			unit = raw[idx:]
			break
		}
	}
	return num, unit
}

func CreateVol(poolName, volName, vtype, capacity string) (*libvirt.StorageVol, error) {
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

	num, unit := parseCapacity(capacity)
	volDesc := &libvirtxml.StorageVolume{
		Type: vtype,
		Name: volName,
		Capacity: &libvirtxml.StorageVolumeSize{
			Unit:  unit,
			Value: num,
		},
		Target: &libvirtxml.StorageVolumeTarget{
			Path: diskPath,
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
	err = vol.Delete(0)
	if err != nil {
		return err
	}
	return vol.Free()
}
