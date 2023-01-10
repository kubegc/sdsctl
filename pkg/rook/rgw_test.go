package rook

import (
	"fmt"
	"testing"
)

func TestGetBucketInfo(t *testing.T) {
	info, _ := GetBucketInfo("ceph-delete-bucket")
	fmt.Printf("%+v\n", info)
}

func TestGetBucketSecret(t *testing.T) {
	GetBucketSecret("ceph-delete-bucket")
}
