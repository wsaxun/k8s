package pkg

import (
	"k8s/pkg/utils"
	"os"
	"path/filepath"
)

type Calico struct {
	utils.K8SSoftware
	DownloadDir string
	Url         string
	PodCIDR     string
}

func (c *Calico) Install() {
	fileName := utils.FileName(c.Url)
	_, err := os.Stat(filepath.Join(c.DownloadDir, fileName))
	var cmd string
	if err != nil {
		cmd = "cd " + c.DownloadDir + " && wget " + c.Url
		utils.Cmd("bash", "-c", cmd)
	}
	cmd = "cd " + c.DownloadDir + " && /usr/bin/cp -f  " + fileName + " calico-private.yml"
	utils.Cmd("bash", "-c", cmd)

	cmd = "cd " + c.DownloadDir + ` &&sed -i 's@# - name: CALICO_IPV4POOL_CIDR@- name: CALICO_IPV4POOL_CIDR@g' calico-private.yml`
	utils.Cmd("bash", "-c", cmd)
	cmd = "cd " + c.DownloadDir + ` &&sed -i 's@#   value: "192.168.0.0/16"@  value: "` + c.PodCIDR + `"@g' calico-private.yml`
	utils.Cmd("bash", "-c", cmd)
	// 旧版本calico.yml 为   value: "192.168.0.0/16"，无 # 号
	cmd = "cd " + c.DownloadDir + ` &&sed -i 's@"192.168.0.0/16"@"` + c.PodCIDR + `"@g' calico-private.yml`
	utils.Cmd("bash", "-c", cmd)

	cmd = "cd " + c.DownloadDir + ` &&sed -i 's@docker.io@quay.io@g' calico-private.yml`
	utils.Cmd("bash", "-c", cmd)

	cmd = "cd " + c.DownloadDir + " && ./kubectl apply -f calico-private.yml"
	utils.Cmd("bash", "-c", cmd)
}
