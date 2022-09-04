package pkg

import (
	"k8s/pkg/utils"
	"os"
	"path/filepath"
)

type Calico struct {
	DownloadDir string
	Url         string
	PodCIDR     string
}

func (c *Calico) InstallCalico() {
	fileName := utils.FileName(c.Url)
	_, err := os.Stat(filepath.Join(c.DownloadDir, fileName))
	var cmd string
	if err != nil {
		cmd = "cd " + c.DownloadDir + " && wget " + c.Url
		utils.Cmd("bash", "-c", cmd)
	}
	cmd = "cd " + c.DownloadDir + ` &&sed -i 's@# - name: CALICO_IPV4POOL_CIDR@- name: CALICO_IPV4POOL_CIDR@g' ` + fileName
	utils.Cmd("bash", "-c", cmd)
	cmd = "cd " + c.DownloadDir + ` &&sed -i 's@#   value: "192.168.0.0/16"@  value: "` + c.PodCIDR + `"@g'` + fileName
	utils.Cmd("bash", "-c", cmd)
	cmd = "cd " + c.DownloadDir + " && ./kubectl apply -f " + fileName
	utils.Cmd("bash", "-c", cmd)
}
