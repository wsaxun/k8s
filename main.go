package main

import (
	"fmt"
	"k8s/pkg/utils"
)

func main() {
	// parser install yml
	data := utils.ParserYml("./configs/install.yml")
	for k, v := range data {
		fmt.Println(k, v)
	}
	// download


	// init env
	//pkg.InitMasterEnv("127.0.0.1")
	//pkg.InitNodeEnv("127.0.0.1")
	//
	//docker := pkg.Docker{
	//	YumRepo:         "http://mirrors.aliyun.com/docker-ce/linux/centos/docker-ce.repo",
	//	DataRoot:        "/var/lib/docker",
	//	RegistryMirrors: "https://mvaav0ar.mirror.aliyuncs.com",
	//}
	//docker.InstallDocker("127.0.0.1")
}
