package pkg

import (
	"github.com/gobuffalo/packr"
	"k8s/pkg/utils"
)

func KubeConfig(downloadCache string, vip string, port int) {
	fileName := "kubeconfig.yml"
	box := packr.NewBox("../template")
	kubeYml, _ := box.FindString(fileName)
	type info struct {
		DownloadDir string
		VIP         string
		Port        int
	}
	content := info{DownloadDir: downloadCache, VIP: vip, Port: port}
	path := utils.Render(content, kubeYml, fileName)
	utils.Playbook(path)
}
