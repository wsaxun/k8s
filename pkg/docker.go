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
	box := packr.NewBox("./template")
	daemonContext, _ := box.FindString("daemon.json")
	utils.Render(d, daemonContext, "daemon.json")

	dockerYml, _ := box.FindString("installDocker.yml")
	content := struct {
		host    string
		YumRepo string
	}{host: host, YumRepo: d.YumRepo}
	utils.Render(content, dockerYml, "installDocker.yml")
}
