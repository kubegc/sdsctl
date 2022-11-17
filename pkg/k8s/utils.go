package k8s

import (
	"context"
	"errors"
	"fmt"
	"github.com/kube-stack/sdsctl/pkg/constant"
	"github.com/kube-stack/sdsctl/pkg/utils"
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

func GetNfsServiceIp() (string, error) {
	client, err := NewClient()
	if err != nil {
		return "", err
	}
	svclist, err := client.CoreV1().Services(constant.RookNamespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return "", err
	}
	for _, svc := range svclist.Items {
		if strings.Contains(svc.Name, "nfs") {
			return svc.Spec.ClusterIP, nil
		}
	}
	return "", errors.New("no nfs service")
}

func CheckNfsMount(nfsSvcIp, path string) bool {
	scmd := fmt.Sprintf("df -h | grep '%s' | awk '{print $6}'", nfsSvcIp)
	cmd := utils.Command{
		Cmd: scmd,
	}
	output, err := cmd.Execute()
	if err != nil || output != "" {
		return strings.Contains(path, output)
	}
	return false
}
