package rook

import (
	"context"
	"fmt"
	"github.com/kube-stack/sdsctl/pkg/constant"
	"github.com/kube-stack/sdsctl/pkg/k8s"
	"github.com/kube-stack/sdsctl/pkg/utils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
)

func CreateOBC(obcName string) error {
	//ksgvr := k8s.NewExternalGvr("objectbucket.io", "v1alpha1", "ObjectBucketClaim")
	//res := make(map[string]interface{})
	//res["generateBucketName"] = "ceph-bkt"
	//res["storageClassName"] = "rook-ceph-delete-bucket"
	//return ksgvr.CreateExternalCrd(context.TODO(), constant.DefaultNamespace, obcName, "spec", res)
	cmd := &utils.Command{
		Cmd: "kubectl apply -f https://gitee.com/syswu/yamls/raw/master/storage/cephrgw/2-object-bucket-claim-delete.yaml",
	}
	_, err := cmd.Execute()
	return err
}

func DeleteOBC() error {
	cmd := &utils.Command{
		Cmd: "kubectl delete -f https://gitee.com/syswu/yamls/raw/master/storage/cephrgw/2-object-bucket-claim-delete.yaml",
	}
	_, err := cmd.Execute()
	return err
}

func CreateBucketStorageClass() error {
	cmd := &utils.Command{
		Cmd: "kubectl apply -f https://gitee.com/syswu/yamls/raw/master/storage/cephrgw/1-storageclass-bucket-delete.yaml",
	}
	_, err := cmd.Execute()
	return err
}

func DeleteBucketStorageClass() error {
	cmd := &utils.Command{
		Cmd: "kubectl delete -f https://gitee.com/syswu/yamls/raw/master/storage/cephrgw/1-storageclass-bucket-delete.yaml",
	}
	_, err := cmd.Execute()
	return err
}

func GetBucketInfo(cmName string) (map[string]string, error) {
	//svcName=$(echo $AWS_HOST | awk -F. '{print $1}') && kubectl get svc -A | grep $svcName | awk '{print $4}'
	//export PORT=$(kubectl -n default get cm ceph-delete-bucket -o jsonpath='{.data.BUCKET_PORT}')
	//export BUCKET_NAME=$(kubectl -n default get cm ceph-delete-bucket -o jsonpath='{.data.BUCKET_NAME}')
	client, err := k8s.NewClient()
	if err != nil {
		return nil, err
	}
	cm, err := client.CoreV1().ConfigMaps(constant.DefaultNamespace).Get(context.TODO(), cmName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	host := cm.Data["BUCKET_HOST"]
	name := cm.Data["BUCKET_NAME"]
	port := cm.Data["BUCKET_PORT"]
	svcName := strings.Split(host, ".")[0]
	scmd := fmt.Sprintf("kubectl get svc -A | grep %s | awk '{print $4}'", svcName)
	cmd := &utils.Command{
		Cmd: scmd,
	}
	ip, _ := cmd.Execute()
	return map[string]string{
		"host": host,
		"name": name,
		"ip":   ip,
		"port": port,
	}, nil
}

func GetBucketSecret(secretName string) (map[string]string, error) {
	//export AWS_ACCESS_KEY_ID=$(kubectl -n default get secret ceph-delete-bucket -o jsonpath='{.data.AWS_ACCESS_KEY_ID}' | base64 --decode)
	//export AWS_SECRET_ACCESS_KEY=$(kubectl -n default get secret ceph-delete-bucket -o jsonpath='{.data.AWS_SECRET_ACCESS_KEY}' | base64 --decode)
	client, err := k8s.NewClient()
	if err != nil {
		return nil, err
	}
	cm, err := client.CoreV1().Secrets(constant.DefaultNamespace).Get(context.TODO(), secretName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	accessId := string(cm.Data["AWS_ACCESS_KEY_ID"])
	accessKey := string(cm.Data["AWS_SECRET_ACCESS_KEY"])
	return map[string]string{
		"access-id":  accessId,
		"access-key": accessKey,
	}, nil
}
