package k8s

import (
	"context"
	"encoding/json"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

var gvr = schema.GroupVersionResource{
	Group:    "doslab.io",
	Version:  "v1",
	Resource: "virtualmachinedisks",
}

type VirtualMachineDiskSpec struct {
	test metav1.Object
}

type VirtualMachineDisk struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              runtime.RawExtension `json:"spec,omitempty"`
}

type VirtualMachineDiskList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []VirtualMachineDisk `json:"items"`
}

func Get(ctx context.Context, client dynamic.Interface, namespace string, name string) (*VirtualMachineDisk, error) {
	utd, err := client.Resource(gvr).Namespace(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	data, err := utd.MarshalJSON()
	fmt.Printf("json:%+v", string(data))
	if err != nil {
		return nil, err
	}
	var vmd VirtualMachineDisk
	if err := json.Unmarshal(data, &vmd); err != nil {
		return nil, err
	}
	return &vmd, nil
}
