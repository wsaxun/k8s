package pkg

import (
	"k8s/pkg/utils"
	"os"
	"path/filepath"
)

type Flannel struct {
	utils.K8SSoftware
	DownloadDir string
	Url         string
	PodCIDR     string
}

func (c *Flannel) Install() {
	fileName := utils.FileName(c.Url)
	_, err := os.Stat(filepath.Join(c.DownloadDir, fileName))
	var cmd string
	if err != nil {
		cmd = "cd " + c.DownloadDir + " && wget " + c.Url
		utils.Cmd("bash", "-c", cmd)
	}
	cmd = "cd " + c.DownloadDir + " && /usr/bin/cp -f  " + fileName + " flannel-private.yml"
	utils.Cmd("bash", "-c", cmd)

	cmd = "cd " + c.DownloadDir + ` &&sed -i 's@"Network": ".*"@"Network": "` + c.PodCIDR + `"@g' flannel-private.yml`
	utils.Cmd("bash", "-c", cmd)

	cmd = "cd " + c.DownloadDir + " && ./kubectl apply -f flannel-private.yml"
	utils.Cmd("bash", "-c", cmd)
}
