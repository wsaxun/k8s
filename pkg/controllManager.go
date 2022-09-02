package pkg

import (
	"github.com/gobuffalo/packr"
	"k8s/pkg/utils"
)

type ControllerManager struct {
	Dir         string
	PodCIDR     string
	DownloadDir string
}

func (c *ControllerManager) InstallControllerManager(host string, inventory string) {
	c.config()
	ymlName := "controllerManager.yml"
	box := packr.NewBox("../template")
	yml, _ := box.FindString(ymlName)

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
	utils.Render(serviceInfo, yml, ymlName)
	//path := utils.Render(content, yml, ymlName)
	//utils.Playbook(path, inventory)
}

func (c *ControllerManager) config() {
	box := packr.NewBox("../template")
	context, _ := box.FindString("softwareConfig/kube-controller-manager.service")
	utils.Render(c, context, "kube-controller-manager.service")
}
