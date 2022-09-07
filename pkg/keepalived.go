package pkg

import (
	"embed"
	"k8s/pkg/utils"
)

type Keepalived struct {
	utils.K8SSoftware
	Interface string
	Host      []string
	Vip       string
}

func (k *Keepalived) Install(host string, inventory string, fs embed.FS) {
	k.config(fs)
	ymlName := "keepalived.yml"
	yml, _ := fs.ReadFile("template/" + ymlName)
	type info struct {
		Host    string
		AllHost []string
	}
	content := info{
		Host:    host,
		AllHost: k.Host,
	}
	path := utils.Render(content, string(yml), ymlName)
	utils.Playbook(path, inventory)
}

func (k *Keepalived) config(fs embed.FS) {
	context, _ := fs.ReadFile("template/softwareConfig/keepalived.conf")
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
		utils.Render(tplData, string(context), "keepalived.conf_"+host)
	}

}
