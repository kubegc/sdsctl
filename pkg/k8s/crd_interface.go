package k8s

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/WANNA959/sdsctl/pkg/constant"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

type KsGvr struct {
	gvr schema.GroupVersionResource
}

type KsCrd struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              runtime.RawExtension `json:"spec,omitempty"`
}

type KsCrdList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []VirtualMachineDisk `json:"items"`
}

func NewKsGvr(crdName string) KsGvr {
	return KsGvr{
		gvr: schema.GroupVersionResource{
			Group:    constant.DefaultGroup,
			Version:  constant.DefaultVersion,
			Resource: crdName,
		},
	}
}

func (ks *KsGvr) Get(ctx context.Context, client dynamic.Interface, namespace string, name string) (*KsCrd, error) {
	utd, err := client.Resource(ks.gvr).Namespace(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	data, err := utd.MarshalJSON()
	fmt.Printf("json:%+v", string(data))
	if err != nil {
		return nil, err
	}
	var kscrd KsCrd
	if err := json.Unmarshal(data, &kscrd); err != nil {
		return nil, err
	}
	return &kscrd, nil
}

func (ks *KsGvr) Exist(ctx context.Context, client dynamic.Interface, namespace string, name string) (bool, error) {
	_, err := ks.Get(ctx, client, namespace, name)
	if err != nil {
		if errors.IsNotFound(err) {
			return false, nil
		} else {
			return false, err
		}
	}
	return true, nil
}

func (ks *KsGvr) Create(ctx context.Context, client dynamic.Interface, namespace string, name string) {

}
