package pkg

import (
	"embed"
	"k8s/pkg/utils"
)

type Kubelet struct {
	utils.K8SSoftware
	CoreDns     string
	Dir         string
	DownloadDir string
}

func (k *Kubelet) Install(host string, inventory string, fs embed.FS) {
	k.config(fs)
	ymlName := "kubelet.yml"
	yml, _ := fs.ReadFile("template/" + ymlName)
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
	path := utils.Render(kubeletInfo, string(yml), ymlName)
	utils.Playbook(path, inventory)
}

func (k *Kubelet) config(fs embed.FS) {
	context, _ := fs.ReadFile("template/softwareConfig/kubelet-config.yml")
	utils.Render(k, string(context), "kubelet-config.yml")
	service, _ := fs.ReadFile("template/softwareConfig/kubelet.service")
	utils.Render(k, string(service), "kubelet.service")
}
