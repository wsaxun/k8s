package pkg

import (
	"github.com/gobuffalo/packr"
	"k8s/pkg/utils"
)

type Keepalived struct {
	Weight    int
	Level     int
	Interface string
	Host      []string
}

func (k *Keepalived) InstallKeepalived(host string, inventory string, config utils.Config) {
	k.config(config)

	//yml, _ := box.FindString(ymlName)
	//type info struct {
	//	Host string
	//}
	//content := info{Host: host}
	//utils.Render(content, yml, ymlName)
	//path := utils.Render(content, yml, ymlName)
	//utils.Playbook(path, inventory)
}

func (k *Keepalived) config(config utils.Config) {
	box := packr.NewBox("../template")
	context, _ := box.FindString("softwareConfig/keepalived.conf")
	type data struct {
		Weight    int
		Level     int
		Interface string
		Host      []string
		LocalHost string
	}

	tplData := data{
		Weight:    2,
		Level:     1,
		Interface: config.Keepalived.Interface,
		Host:      nil,
		LocalHost: nil,
	}
	for index, host := range k.Host {
		tplData.Level += index
		tplData.LocalHost = host
		tplData.Host = k.Host
		utils.Render(tplData, context, "keepalived.conf"+host)
	}

}
