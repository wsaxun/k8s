package pkg

import (
	"github.com/gobuffalo/packr"
	"k8s/pkg/utils"
)

type Kubelet struct {
	utils.K8SSoftware
	CoreDns     string
	Dir         string
	DownloadDir string
}

func (k *Kubelet) Install(host string, inventory string) {
	k.config()
	ymlName := "kubelet.yml"
	box := packr.NewBox("../template")
	yml, _ := box.FindString(ymlName)
	type info struct {
		CoreDns     string
		Dir         string
		DownloadDir string
		Host        string
	}
	kubeletInfo := info{
		CoreDns:     k.CoreDns,
		Dir:         k.Dir,
		DownloadDir: k.DownloadDir,
		Host:        host,
	}
	utils.Render(kubeletInfo, yml, ymlName)
	//path := utils.Render(content, yml, ymlName)
	//utils.Playbook(path, inventory)
}

func (k *Kubelet) config() {
	box := packr.NewBox("../template")
	context, _ := box.FindString("softwareConfig/kubelet-config.yml")
	utils.Render(k, context, "kubelet-config.yml")
	service, _ := box.FindString("softwareConfig/kubelet.service")
	utils.Render(k, service, "kubelet.service")
}
