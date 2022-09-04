package pkg

import (
	"fmt"
	"github.com/gobuffalo/packr"
	"k8s/pkg/utils"
	"strconv"
)

type Etcd struct {
	utils.K8SSoftware
	DataDir     string
	Host        []string
	Dir         string
	DownloadDir string
	EtcdName    string
}

func (e *Etcd) Install(host string, inventory string) {
	e.config()
	ymlName := "etcd.yml"
	box := packr.NewBox("../template")
	yml, _ := box.FindString(ymlName)
	type info struct {
		Host        string
		AllHost     []string
		DownloadDir string
		Dir         string
		EtcdName    string
	}
	content := info{
		Host:        host,
		AllHost:     e.Host,
		Dir:         e.Dir,
		DownloadDir: e.DownloadDir,
		EtcdName:    e.EtcdName,
	}
	utils.Render(content, yml, ymlName)
	//path := utils.Render(content, yml, ymlName)
	//utils.Playbook(path, inventory)
}

func (e *Etcd) config() {
	box := packr.NewBox("../template")
	context, _ := box.FindString("softwareConfig/etcd.config.yml")
	type data struct {
		Name      string
		DataDir   string
		LocalHost string
		Cluster   string
	}

	tplData := data{
		Name:      "",
		DataDir:   e.DataDir,
		LocalHost: "",
		Cluster:   "",
	}
	// k8s-master01=https://192.168.58.129:2380,k8s-master02=https://192.168.58.130:2380,k8s-master03=https://192.168.58.131:2380
	cluster := ""
	for index, host := range e.Host {
		if cluster == "" {
			cluster = cluster + fmt.Sprintf("etcd%d=https://%s:2380", index, host)
		} else {
			cluster = cluster + fmt.Sprintf(",etcd%d=https://%s:2380", index, host)
		}
	}
	tplData.Cluster = cluster
	for index, host := range e.Host {
		tplData.Name = "etcd" + strconv.Itoa(index)
		tplData.LocalHost = host
		utils.Render(tplData, context, "etcd.config.yml_"+host)
	}

	service, _ := box.FindString("softwareConfig/etcd.service")
	type info struct {
		Dir string
	}
	serviceInfo := info{Dir: e.Dir}
	utils.Render(serviceInfo, service, "etcd.service")
}
