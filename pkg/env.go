package pkg

import (
	"github.com/gobuffalo/packr"
	"k8s/pkg/utils"
)

func initEnv(host string, yml string) {
	type info struct {
		Host string
	}

	hosts := info{Host: host}

	box := packr.NewBox("../template")
	initMasterYml, _ := box.FindString(yml)
	path := utils.Render(hosts, initMasterYml, yml)
	utils.Playbook(path)
}

func InitMasterEnv(host string) {
	initEnv(host, "initMaster.yml")
}

func InitNodeEnv(host string) {
	initEnv(host, "initNode.yml")
}
