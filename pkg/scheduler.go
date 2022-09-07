package pkg

import (
	"embed"
	"k8s/pkg/utils"
)

type Scheduler struct {
	utils.K8SSoftware
	Dir         string
	DownloadDir string
}

func (s *Scheduler) Install(host string, inventory string, fs embed.FS) {
	s.config(fs)
	ymlName := "scheduler.yml"
	yml, _ := fs.ReadFile("template/" + ymlName)

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
	path := utils.Render(serviceInfo, string(yml), ymlName)
	utils.Playbook(path, inventory)
}

func (s *Scheduler) config(fs embed.FS) {
	context, _ := fs.ReadFile("template/softwareConfig/kube-scheduler.service")
	utils.Render(s, string(context), "kube-scheduler.service")
}
