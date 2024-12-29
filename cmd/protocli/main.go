package main

import (
	"flag"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	protoPath := flag.String("path", "proto", "proto_path proto文件所在目录")
	goOutPath := flag.String("out", "/Users/liujinglong/project/msproject/project-grpc", "go_out pb.go文件输出目录")
	goPackage := flag.String("package", "github.com/msprojectlb/project-grpc", "go_package 微服务项目地址")
	flag.Parse()
	_, err := exec.LookPath("protoc")
	if err != nil {
		log.Fatal("protoc is not installed")
	}
	_, err = exec.LookPath("protoc-gen-go")
	if err != nil {
		log.Fatal("protoc-gen-go is not installed")
	}
	_, err = exec.LookPath("protoc-gen-go-grpc")
	if err != nil {
		log.Fatal("protoc-gen-go-grpc is not installed")
	}
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	protoFilePath := filepath.Join(pwd, *protoPath)
	args := []string{
		"--proto_path=" + protoFilePath,
		"--go_out=" + *goOutPath,
		"--go-grpc_out=" + *goOutPath,
		"--go_opt=module=" + *goPackage,
		"--go-grpc_opt=module=" + *goPackage,
	}
	if _, err := os.Stat(protoFilePath); os.IsNotExist(err) {
		log.Fatal("filepath not exist " + protoFilePath)
	}
	err = filepath.Walk(protoFilePath, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && filepath.Ext(path) == ".proto" {
			args = append(args, path)
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	cmd := exec.Command("protoc", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}
