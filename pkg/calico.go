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
	_, err := os.Stat(filepath.Join(c.DownloadDir, "calico.yaml"))
	if err != nil {
		cmd := "cd " + c.DownloadDir + " && curl " + c.Url + " -o calico.yaml"
		utils.Cmd("bash", "-c", cmd)
	}
	cmd = "cd " + c.DownloadDir + ` &&sed -i 's@# - name: CALICO_IPV4POOL_CIDR@- name: CALICO_IPV4POOL_CIDR@g' calico.yaml`
	utils.Cmd("bash", "-c", cmd)
	cmd = "cd " + c.DownloadDir + ` &&sed -i 's@#   value: "192.168.0.0/16"@  value: "` + c.PodCIDR + `"@g' calico.yaml`
	utils.Cmd("bash", "-c", cmd)
	cmd = "cd " + c.DownloadDir + " && ./kubectl apply -f calico.yml"
	utils.Cmd("bash", "-c", cmd)
}
