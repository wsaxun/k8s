package pkg

import (
	"embed"
	"k8s/pkg/utils"
)

type ControllerManager struct {
	utils.K8SSoftware
	Dir         string
	PodCIDR     string
	DownloadDir string
}

func (c *ControllerManager) Install(host string, inventory string, fs embed.FS) {
	c.config(fs)
	ymlName := "controllerManager.yml"
	yml, _ := fs.ReadFile("template/controllerManager.yml")

	type info struct {
		Dir         string
		Host        string
		DownloadDir string
	}
	serviceInfo := info{
		Dir:         c.Dir,
		Host:        host,
		DownloadDir: c.DownloadDir,
	}
	path := utils.Render(serviceInfo, string(yml), ymlName)
	utils.Playbook(path, inventory)
}

func (c *ControllerManager) config(fs embed.FS) {
	context, _ := fs.ReadFile("template/softwareConfig/kube-controller-manager.service")
	utils.Render(c, string(context), "kube-controller-manager.service")
}
