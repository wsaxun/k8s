package pkg

import (
	"github.com/gobuffalo/packr"
	"k8s/pkg/utils"
)

type Docker struct {
	YumRepo         string
	DataRoot        string
	RegistryMirrors string
}

func (d *Docker) InstallDocker(host string) {
	templateName := "daemon.json"
	ymlName := "installDocker.yml"
	box := packr.NewBox("../template")
	daemonContext, _ := box.FindString(templateName)
	utils.Render(d, daemonContext, templateName)

	dockerYml, _ := box.FindString(ymlName)
	type info struct {
		Host    string
		YumRepo string
	}
	content := info{Host: host, YumRepo: d.YumRepo}
	path := utils.Render(content, dockerYml, ymlName)
	utils.Playbook(path)
}
