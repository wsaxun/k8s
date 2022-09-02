package pkg

import (
	"github.com/gobuffalo/packr"
	"k8s/pkg/utils"
)

type ApiServer struct {
	Host        []string
	Dir         string
	DownloadDir string
	ServiceCIDR string
	EtcdHost    []string
}

func (a *ApiServer) InstallApiServer(host string, inventory string) {
	a.config()
	ymlName := "apiServer.yml"
	box := packr.NewBox("../template")
	yml, _ := box.FindString(ymlName)

	type info struct {
		Dir         string
		Host        string
		DownloadDir string
		AllHost     []string
	}
	serviceInfo := info{
		Dir:         a.Dir,
		Host:        host,
		DownloadDir: a.DownloadDir,
		AllHost:     a.Host,
	}
	utils.Render(serviceInfo, yml, ymlName)
	//path := utils.Render(content, yml, ymlName)
	//utils.Playbook(path, inventory)
}

func (a *ApiServer) config() {
	box := packr.NewBox("../template")
	context, _ := box.FindString("softwareConfig/kube-apiserver.service")
	type data struct {
		Dir         string
		LocalHost   string
		ServiceCIDR string
		EtcdHost    string
	}

	var etcdUrl string
	for _, host := range a.EtcdHost {
		if etcdUrl == "" {
			etcdUrl = etcdUrl + "https://" + host + ":2379"
		} else {
			etcdUrl = etcdUrl + ",https://" + host + ":2379"
		}
	}

	tplData := data{
		Dir:         a.Dir,
		LocalHost:   "",
		ServiceCIDR: a.ServiceCIDR,
		EtcdHost:    etcdUrl,
	}
	for _, host := range a.Host {
		tplData.LocalHost = host
		utils.Render(tplData, context, "kube-apiserver.service_"+host)
	}
}
