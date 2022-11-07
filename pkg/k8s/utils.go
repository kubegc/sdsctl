package k8s

import (
	"context"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
	"strings"
)

func GetVMHostName() string {
	name, _ := os.Hostname()
	return fmt.Sprintf("vm.%s", strings.ToLower(name))
}

func GetIPByNodeName(nodeName string) (string, error) {
	client, err := NewClient()
	if err != nil {
		return "", err
	}
	nodeInfo, err := client.CoreV1().Nodes().Get(context.TODO(), nodeName, metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	annotations := nodeInfo.GetObjectMeta().GetAnnotations()
	return annotations["THISIP"], nil
}