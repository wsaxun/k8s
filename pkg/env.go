package pkg

import (
	"github.com/gobuffalo/packr"
	"k8s/pkg/utils"
)

func initEnv(host string, yml string, inventory string) {
	type info struct {
		Host string
	}

	hosts := info{Host: host}

	box := packr.NewBox("../template")
	initMasterYml, _ := box.FindString(yml)
	utils.Render(hosts, initMasterYml, yml)
	//path := utils.Render(hosts, initMasterYml, yml)
	//utils.Playbook(path, inventory)
}

func InitMasterEnv(host string, inventory string) {
	initEnv(host, "initMaster.yml", inventory)
}

func InitNodeEnv(host string, inventory string) {
	initEnv(host, "initNode.yml", inventory)
}
