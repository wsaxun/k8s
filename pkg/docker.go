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
	ymlName := "installDocker.yml"
	box := packr.NewBox("../template")
	daemonContext, _ := box.FindString("softwareConfig/daemon.json")
	utils.Render(d, daemonContext, "daemon.json")

	dockerYml, _ := box.FindString(ymlName)
	type info struct {
		Host    string
		YumRepo string
	}
	content := info{Host: host, YumRepo: d.YumRepo}
	path := utils.Render(content, dockerYml, ymlName)
	utils.Playbook(path)
}
