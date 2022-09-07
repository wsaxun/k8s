package pkg

import (
	"embed"
	"fmt"
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
	Record      map[string]string
}

func (e *Etcd) Install(host string, inventory string, fs embed.FS) {
	e.config(fs)

	//add hosts
	e.addEtcdHosts(host, inventory, fs)

	ymlName := "etcd.yml"
	yml, _ := fs.ReadFile("template/etcd.yml")
	type info struct {
		Host        string
		AllHost     []string
		DownloadDir string
		Dir         string
		EtcdName    string
		DataDir     string
	}
	content := info{
		Host:        host,
		AllHost:     e.Host,
		Dir:         e.Dir,
		DownloadDir: e.DownloadDir,
		EtcdName:    e.EtcdName,
		DataDir:     e.DataDir,
	}
	path := utils.Render(content, string(yml), ymlName)
	utils.Playbook(path, inventory)
}

func (e *Etcd) config(fs embed.FS) {
	context, _ := fs.ReadFile("template/softwareConfig/etcd.config.yml")
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
	recordMap := make(map[string]string)
	for index, host := range e.Host {
		if cluster == "" {
			cluster = cluster + fmt.Sprintf("etcd%d=https://%s:2380", index, host)
		} else {
			cluster = cluster + fmt.Sprintf(",etcd%d=https://%s:2380", index, host)
		}
		recordMap[fmt.Sprintf("etcd%d", index)] = host
	}
	e.Record = recordMap
	tplData.Cluster = cluster
	for index, host := range e.Host {
		tplData.Name = "etcd" + strconv.Itoa(index)
		tplData.LocalHost = host
		utils.Render(tplData, string(context), "etcd.config.yml_"+host)
	}

	service, _ := fs.ReadFile("template/softwareConfig/etcd.service")
	type info struct {
		Dir string
	}
	serviceInfo := info{Dir: e.Dir}
	utils.Render(serviceInfo, string(service), "etcd.service")
}

func (e *Etcd) addEtcdHosts(host, inventory string, fs embed.FS) {
	ymlName := "addEtcdHosts.yml"
	yml, _ := fs.ReadFile("template/" + ymlName)
	type info struct {
		Host   string
		Record map[string]string
	}
	content := info{
		Host:   host,
		Record: e.Record,
	}
	path := utils.Render(content, string(yml), ymlName)
	utils.Playbook(path, inventory)
}
