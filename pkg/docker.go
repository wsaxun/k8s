package pkg

import (
	"fmt"
	"k8s/pkg/utils"
)

type Docker struct {
	YumRepo         string
	DataRoot        string
	RegistryMirrors string
}

// config yum repo
func (d *Docker) ConfigYum(host string) {
	utils.Exec(host, "shell", "yum-config-manager --add-repo "+d.YumRepo)
	utils.Exec(host, "shell", "yum makecache fast")
}

func (d *Docker) InstallDocker(host string) {
	utils.Exec(host, "yum", "name=docker-ce state=installed")
	utils.Exec(host, "file", "path=/etc/docker state=directory")
	utils.Exec(host, "file", "path=/etc/docke/daemon.json state=touch")
	conf := fmt.Sprintf("{\"registry-mirrors\":[\"%s\"],\"exec-opts\":[\"native.cgroupdriver=systemd\"],\"data-root\":\"%s\",\"log-driver\":\"json-file\",\"log-opts\":{\"max-size\":\"2048m\",\"max-file\":\"5\"}}", d.RegistryMirrors, d.DataRoot)
	utils.Exec(host, "copy", "content="+conf)
}

func (d *Docker) StartDocker(host string) {
	utils.Exec(host, "service", "name=docker state=started enable=yes")
}
