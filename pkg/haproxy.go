package pkg

import (
	"github.com/gobuffalo/packr"
	"k8s/pkg/utils"
)

type Haproxy struct {
	utils.K8SSoftware
	Port     int
	HostInfo map[string]string
}

func (h *Haproxy) Install(host string, inventory string) {
	ymlName := "haproxy.yml"
	box := packr.NewBox("../template")
	context, _ := box.FindString("softwareConfig/haproxy.cfg")
	utils.Render(h, context, "haproxy.cfg")

	yml, _ := box.FindString(ymlName)
	type info struct {
		Host string
	}
	content := info{Host: host}
	utils.Render(content, yml, ymlName)
	//path := utils.Render(content, yml, ymlName)
	//utils.Playbook(path, inventory)
}
