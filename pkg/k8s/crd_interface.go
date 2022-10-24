package k8s

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/WANNA959/sdsctl/pkg/constant"
	"github.com/tidwall/sjson"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"path/filepath"
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
	Items           []KsCrd `json:"items"`
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

var Client dynamic.Interface

func GetClient() (dynamic.Interface, error) {
	var err error
	if Client == nil {
		Client, err = NewClient()
		if err != nil {
			return nil, err
		}
	}
	return Client, nil
}

func NewClient() (dynamic.Interface, error) {
	kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}

	client, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (ks *KsGvr) Get(ctx context.Context, namespace string, name string) (*KsCrd, error) {
	client, err := GetClient()
	if err != nil {
		return nil, err
	}
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

func (ks *KsGvr) Exist(ctx context.Context, namespace string, name string) (bool, error) {
	_, err := ks.Get(ctx, namespace, name)
	if err != nil {
		if errors.IsNotFound(err) {
			return false, nil
		} else {
			return false, err
		}
	}
	return true, nil
}

func (ks *KsGvr) Update(ctx context.Context, namespace, name, key string, value interface{}) error {
	client, err := GetClient()
	if err != nil {
		return err
	}

	// get old crd
	utd, err := client.Resource(ks.gvr).Namespace(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return err
	}

	obj := &unstructured.Unstructured{}
	obj.SetResourceVersion(utd.GetResourceVersion())

	data, err := utd.MarshalJSON()
	if err != nil {
		return err
	}
	var kscrd KsCrd
	err = json.Unmarshal(data, &kscrd)
	if err != nil {
		return err
	}

	// update spec
	//fmt.Printf("before:%s\n", string(kscrd.Spec.Raw))
	bytes, err := sjson.SetBytes(kscrd.Spec.Raw, key, value)
	if err != nil {
		return err
	}
	kscrd.Spec.Raw = bytes
	//fmt.Printf("after:%s\n", string(kscrd.Spec.Raw))
	// docode for bytes
	marshal, err := json.Marshal(kscrd)
	decoder := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
	if _, _, err = decoder.Decode(marshal, nil, obj); err != nil {
		return err
	}
	// write back to k8s
	_, err = client.Resource(ks.gvr).Namespace(namespace).Update(ctx, obj, metav1.UpdateOptions{})
	if err != nil {
		return err
	}
	return nil

}

func (ks *KsGvr) Delete(ctx context.Context, namespace string, name string) error {
	client, err := GetClient()
	if err != nil {
		return err
	}
	return client.Resource(ks.gvr).Namespace(namespace).Delete(ctx, name, metav1.DeleteOptions{})
}
