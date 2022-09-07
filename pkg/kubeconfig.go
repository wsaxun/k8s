package pkg

import (
	"embed"
	"k8s/pkg/utils"
)

func KubeConfig(downloadCache string, vip string, port int, inventory string, fs embed.FS) {
	fileName := "kubeconfig.yml"
	kubeYml, _ := fs.ReadFile("template/" + fileName)
	type info struct {
		DownloadDir string
		VIP         string
		Port        int
	}
	content := info{DownloadDir: downloadCache, VIP: vip, Port: port}
	path := utils.Render(content, string(kubeYml), fileName)
	utils.Playbook(path, inventory)
}
