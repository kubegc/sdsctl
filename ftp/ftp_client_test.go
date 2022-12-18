package ftp

import (
	"fmt"
	"testing"
)

func TestFtpClient_ListDir(t *testing.T) {
	client, _ := NewFtpClient("192.168.100.100", "21", "ftpuser", "ftpuser")
	dir, err := client.ListDir("/")
	fmt.Printf("dir:%+v\n", dir)
	fmt.Printf("err:%+v\n", err)
}

func TestFtpClient_IsDirExsit(t *testing.T) {
	client, _ := NewFtpClient("192.168.100.100", "21", "ftpuser", "ftpuser")
	exsit := client.IsDirExsit("/test")
	fmt.Printf("exsit:%+v\n", exsit)
}

func TestFtpClient_IsFileExsit(t *testing.T) {
	client, _ := NewFtpClient("192.168.100.100", "21", "ftpuser", "ftpuser")
	exsit := client.IsFileExsit("/name.txt")
	fmt.Printf("exsit:%+v\n", exsit)
}

func TestFtpClient_Mkdir(t *testing.T) {
	client, _ := NewFtpClient("192.168.100.100", "21", "ftpuser", "ftpuser")
	err := client.Mkdir("/test2")
	fmt.Printf("err:%+v\n", err)
}

func TestFtpClient_Rename(t *testing.T) {
	client, err := NewFtpClient("192.168.100.100", "21", "ftpuser", "ftpuser")
	//fmt.Printf("err:%+v\n", err)
	fmt.Printf("client:%+v\n", client)
	err = client.Rename("/", "join.sh", "back.sh")
	fmt.Printf("err:%+v\n", err)
}

func TestFtpClient_UploadFile(t *testing.T) {
	client, err := NewFtpClient("192.168.100.100", "21", "ftpuser", "ftpuser")
	//fmt.Printf("err:%+v\n", err)
	fmt.Printf("client:%+v\n", client)
	err = client.UploadFile("/root/join.sh", "/")
	fmt.Printf("err:%+v\n", err)
}

func TestFtpClient_DownloadFile(t *testing.T) {
	client, err := NewFtpClient("192.168.100.100", "21", "ftpuser", "ftpuser")
	//fmt.Printf("err:%+v\n", err)
	fmt.Printf("client:%+v\n", client)
	err = client.DownloadFile("/root", "/back.sh")
	fmt.Printf("err:%+v\n", err)
}
