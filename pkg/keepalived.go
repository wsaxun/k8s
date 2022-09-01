package pkg

import (
	"github.com/gobuffalo/packr"
	"k8s/pkg/utils"
)

type Keepalived struct {
	Interface string
	Host      []string
}

func (k *Keepalived) InstallKeepalived(host string, inventory string) {
	k.config()
	ymlName := "keepalived.yml"
	yml, _ := box.FindString(ymlName)
	type info struct {
		Host    string
		AllHost []string
	}
	content := info{
		Host:    host,
		AllHost: k.Host,
	}
	utils.Render(content, yml, ymlName)
	//path := utils.Render(content, yml, ymlName)
	//utils.Playbook(path, inventory)
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
	}

	tplData := data{
		Weight:    2,
		Level:     1,
		Interface: k.Interface,
		LocalHost: "",
		HostInfo:  k.Host,
	}
	for index, host := range k.Host {
		tplData.Level += index
		tplData.LocalHost = host
		utils.Render(tplData, context, "keepalived.conf_"+host)
	}

}
