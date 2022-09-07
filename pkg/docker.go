package pkg

import (
	"embed"
	"k8s/pkg/utils"
)

type Docker struct {
	utils.K8SSoftware
	YumRepo         string
	DataRoot        string
	RegistryMirrors string
}

func (d *Docker) Install(host string, inventory string, fs embed.FS) {
	ymlName := "installDocker.yml"
	daemonContext, _ := fs.ReadFile("template/softwareConfig/daemon.json")
	utils.Render(d, string(daemonContext), "daemon.json")

	dockerYml, _ := fs.ReadFile("template/" + ymlName)
	type info struct {
		Host    string
		YumRepo string
	}
	content := info{Host: host, YumRepo: d.YumRepo}
	path := utils.Render(content, string(dockerYml), ymlName)
	utils.Playbook(path, inventory)
}
