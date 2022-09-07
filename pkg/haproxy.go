package pkg

import (
	"embed"
	"k8s/pkg/utils"
)

type Haproxy struct {
	utils.K8SSoftware
	Port     int
	HostInfo map[string]string
}

func (h *Haproxy) Install(host string, inventory string, fs embed.FS) {
	ymlName := "haproxy.yml"
	context, _ := fs.ReadFile("template/softwareConfig/haproxy.cfg")
	utils.Render(h, string(context), "haproxy.cfg")

	yml, _ := fs.ReadFile("template/" + ymlName)
	type info struct {
		Host string
	}
	content := info{Host: host}
	path := utils.Render(content, string(yml), ymlName)
	utils.Playbook(path, inventory)
}
