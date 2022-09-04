package pkg

import (
	"github.com/gobuffalo/packr"
	"k8s/pkg/utils"
)

type Keepalived struct {
	utils.K8SSoftware
	Interface string
	Host      []string
	Vip       string
}

func (k *Keepalived) Install(host string, inventory string) {
	k.config()
	ymlName := "keepalived.yml"
	box := packr.NewBox("../template")
	yml, _ := box.FindString(ymlName)
	type info struct {
		Host    string
		AllHost []string
	}
	content := info{
		Host:    host,
		AllHost: k.Host,
	}
	path := utils.Render(content, yml, ymlName)
	utils.Playbook(path, inventory)
}

func (k *Keepalived) config() {
	box := packr.NewBox("../template")
	context, _ := box.FindString("softwareConfig/keepalived.conf")
	type data struct {
		Weight    int
		Level     int
		Interface string
		HostInfo  []string
		LocalHost string
		Vip       string
	}

	tplData := data{
		Weight:    2,
		Level:     1,
		Interface: k.Interface,
		LocalHost: "",
		HostInfo:  k.Host,
		Vip:       k.Vip,
	}
	for _, host := range k.Host {
		tplData.Level += 1
		tplData.LocalHost = host
		utils.Render(tplData, context, "keepalived.conf_"+host)
	}

}
