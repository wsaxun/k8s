package pkg

import (
	"embed"
	"k8s/pkg/utils"
)

func initEnv(host string, yml string, inventory string, fs embed.FS) {
	type info struct {
		Host string
	}

	hosts := info{Host: host}

	initMasterYml, _ := fs.ReadFile("template/" + yml)
	path := utils.Render(hosts, string(initMasterYml), yml)
	utils.Playbook(path, inventory)
}

func InitMasterEnv(host string, inventory string, fs embed.FS) {
	initEnv(host, "initMaster.yml", inventory, fs)
}

func InitNodeEnv(host string, inventory string, fs embed.FS) {
	initEnv(host, "initNode.yml", inventory, fs)
}
