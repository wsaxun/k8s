package pkg

import (
	"embed"
	"k8s/pkg/utils"
)

type Containerd struct {
	utils.K8SSoftware
	YumRepo         string
	DataRoot        string
	RegistryMirrors string
	SandboxImage    string
}

func (d *Containerd) Install(host string, inventory string, fs embed.FS) {
	ymlName := "installContainerd.yml"
	contanierdYml, _ := fs.ReadFile("template/" + ymlName)
	type info struct {
		Host            string
		YumRepo         string
		DataRoot        string
		RegistryMirrors string
		SandboxImage    string
	}
	content := info{
		Host:            host,
		YumRepo:         d.YumRepo,
		DataRoot:        d.DataRoot,
		RegistryMirrors: d.RegistryMirrors,
		SandboxImage:    d.SandboxImage,
	}
	path := utils.Render(content, string(contanierdYml), ymlName)
	utils.Playbook(path, inventory)
}
