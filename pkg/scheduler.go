package pkg

import (
	"github.com/gobuffalo/packr"
	"k8s/pkg/utils"
)

type Scheduler struct {
	utils.K8SSoftware
	Dir         string
	DownloadDir string
}

func (s *Scheduler) Install(host string, inventory string) {
	s.config()
	ymlName := "scheduler.yml"
	box := packr.NewBox("../template")
	yml, _ := box.FindString(ymlName)

	type info struct {
		Dir         string
		Host        string
		DownloadDir string
	}
	serviceInfo := info{
		Dir:         s.Dir,
		Host:        host,
		DownloadDir: s.DownloadDir,
	}
	utils.Render(serviceInfo, yml, ymlName)
	//path := utils.Render(content, yml, ymlName)
	//utils.Playbook(path, inventory)
}

func (s *Scheduler) config() {
	box := packr.NewBox("../template")
	context, _ := box.FindString("softwareConfig/kube-scheduler.service")
	utils.Render(s, context, "kube-scheduler.service")
}
