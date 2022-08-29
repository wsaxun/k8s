package main

import (
	"fmt"
	"k8s/pkg"
	"k8s/pkg/utils"
)

func main() {
	data := utils.ParserYml("./configs/install.yml")
	for k, v := range data {
		fmt.Println(k, v)
	}
	err := utils.Exec("127.0.0.1", "command", "ping 127.0.0.1 -c 2")
	fmt.Println(err)

	pkg.Docker{
		YumRepo:         "http://mirrors.aliyun.com/docker-ce/linux/centos/docker-ce.repo",
		DataRoot:        "/var/lib/docker",
		RegistryMirrors: "https://mvaav0ar.mirror.aliyuncs.com",
	}.InstallDocker("127.0.0.1")
}
