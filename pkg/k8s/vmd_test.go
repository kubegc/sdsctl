package k8s

import (
	"context"
	"fmt"
	"k8s.io/client-go/dynamic"

	"github.com/tidwall/gjson"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"path/filepath"
	"testing"
)

func TestGet(t *testing.T) {
	kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err)
	}

	client, err := dynamic.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	vmd, err := Get(context.TODO(), client, "default", "disktest131")
	if err != nil {
		panic(err)
	}
	fmt.Printf("disktest131: %+v\n", vmd)
	fmt.Printf("disktest131 spec: %+v\n", string(vmd.Spec.Raw))

	parse := gjson.ParseBytes(vmd.Spec.Raw)
	nodeName := parse.Get("nodeName")
	msg := parse.Get("status.conditions.state.waiting.message")
	fmt.Printf("disktest131 spec nodename: %+v\n", nodeName)
	fmt.Printf("disktest131 spec msg: %+v\n", msg)
}
