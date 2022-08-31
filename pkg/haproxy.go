package pkg

import (
	"github.com/gobuffalo/packr"
	"k8s/pkg/utils"
)

type Haproxy struct {
	Port     int
	HostInfo map[string]string
}

func (h *Haproxy) InstallHaproxy(host string) {
	//ymlName := "haproxy.yml"
	box := packr.NewBox("../template")
	context, _ := box.FindString("softwareConfig/haproxy.cfg")
	utils.Render(h, context, "haproxy.cfg")

	//yml, _ := box.FindString(ymlName)
	//type info struct {
	//	Host    string
	//	YumRepo string
	//}
	//content := info{Host: host, YumRepo: d.YumRepo}
	//path := utils.Render(content, yml, ymlName)
	//utils.Playbook(path)
}
