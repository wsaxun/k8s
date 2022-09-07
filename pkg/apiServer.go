package pkg

import (
	"embed"
	"k8s/pkg/utils"
)

type ApiServer struct {
	utils.K8SSoftware
	Host        []string
	Dir         string
	DownloadDir string
	ServiceCIDR string
	EtcdHost    []string
}

func (a *ApiServer) Install(host string, inventory string, fs embed.FS) {
	a.config(fs)
	ymlName := "apiServer.yml"
	yml, _ := fs.ReadFile("template/apiServer.yml")

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
	path := utils.Render(serviceInfo, string(yml), ymlName)
	utils.Playbook(path, inventory)
}

func (a *ApiServer) config(fs embed.FS) {
	context, _ := fs.ReadFile("template/softwareConfig/kube-apiserver.service")
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
		utils.Render(tplData, string(context), "kube-apiserver.service_"+host)
	}
}
