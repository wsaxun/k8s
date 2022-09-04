package pkg

import (
	"github.com/gobuffalo/packr"
	"k8s/pkg/utils"
)

type CoreDns struct {
	utils.K8SSoftware
	Dns         string
	DownloadDir string
}

func (c *CoreDns) Install() {
	box := packr.NewBox("../template")
	context, _ := box.FindString("softwareConfig/coredns.yml")
	path := utils.Render(c, context, "coredns.yml")
	cmd := "cd " + c.DownloadDir + " && ./kubectl apply -f  " + path
	utils.Cmd("bash", "-c", cmd)
}
