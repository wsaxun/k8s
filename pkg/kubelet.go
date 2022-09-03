package pkg

import (
	"github.com/gobuffalo/packr"
	"k8s/pkg/utils"
)

type Kubelet struct {
	CoreDns     string
	Dir         string
	DownloadDir string
}

func (k *Kubelet) InstallKubelet(host string, inventory string) {
	k.config()
	ymlName := "kubelet.yml"
	box := packr.NewBox("../template")
	yml, _ := box.FindString(ymlName)

	utils.Render(k, yml, ymlName)
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
