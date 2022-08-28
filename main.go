package main

import (
	"fmt"
	"github.com/gobuffalo/packr"
	"k8s/pkg/utils"
)

func main() {
	box := packr.NewBox("./template")
	daemonTemplate, _ := box.Open("daemon.json")
	defer daemonTemplate.Close()
	fmt.Println(daemonTemplate)

	data := utils.ParserYml("./configs/install.yml")
	for k, v := range data {
		fmt.Println(k, v)
	}
	err := utils.Exec("127.0.0.1", "command", "ping 127.0.0.1 -c 2")
	fmt.Println(err)
}
