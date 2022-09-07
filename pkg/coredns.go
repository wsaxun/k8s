package pkg

import (
	"embed"
	"k8s/pkg/utils"
)

type CoreDns struct {
	utils.K8SSoftware
	Dns         string
	DownloadDir string
}

func (c *CoreDns) Install(fs embed.FS) {
	context, _ := fs.ReadFile("template/softwareConfig/coredns.yml")
	path := utils.Render(c, string(context), "coredns.yml")
	cmd := "cd " + c.DownloadDir + " && ./kubectl apply -f  " + path
	utils.Cmd("bash", "-c", cmd)
}
