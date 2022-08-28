package main

import (
	"fmt"
	"github.com/gobuffalo/packr"
	"k8s/pkg/utils"
)

func main() {
	packr.NewBox("./template")
	data := utils.ParserYml("E:\\code\\go\\k8s\\configs\\install.yml")
	for k, v := range data {
		fmt.Println(k, v)
	}
	conf := fmt.Sprintf("{\"registry-mirrors\":[\"%s\"],\"exec-opts\":[\"native.cgroupdriver=systemd\"],\"data-root\":\"%s\",\"log-driver\":\"json-file\",\"log-opts\":{\"max-size\":\"2048m\",\"max-file\":\"5\"}}", "https://xx", "/data/docker")
	fmt.Println(conf)
	err := utils.Exec("127.0.0.1", "command", "ping 127.0.0.1 -c 2")
	fmt.Println(err)
}
