package rook

import (
	"context"
	"github.com/kube-stack/sdsctl/pkg/constant"
	"github.com/kube-stack/sdsctl/pkg/k8s"
)

func CreateRbdPool(poolName string) error {
	ksgvr := k8s.NewExternalGvr(constant.DefaultRookGroup, constant.DefaultRookVersion, constant.CephBlockPoolS_Kinds)
	res := make(map[string]interface{})
	res["failureDomain"] = "host"
	res["replicated"] = map[string]interface{}{
		"size":                   3,
		"requireSafeReplicaSize": true,
	}
	return ksgvr.CreateExternalCrd(context.TODO(), constant.RookNamespace, poolName, "spec", res)
}
