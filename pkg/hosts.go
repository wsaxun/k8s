package pkg

import (
	"github.com/gobuffalo/packr"
	"k8s/pkg/utils"
)

func Inventory(config utils.Config) string {
	type info struct {
		Master                []string
		EtcdHost              []string
		ApiServerHost         []string
		SchedulerHost         []string
		ControllerManagerHost []string
		NodeHost              []string
		HaproxyHost           []string
		KeepalivedHost        []string
	}

	var allHost []string
	var master []string
	var etcdHost []string
	var apiServerHost []string
	var schedulerHost []string
	var controllerManagerHost []string
	nodeHost := config.K8s.Node.Hosts
	haproxyHost := config.Haproxy.Hosts
	keepalivedHost := config.Keepalived.Hosts
	for _, v := range config.K8s.Master.Components {
		switch v.Name {
		case "etcd":
			etcdHost = v.Hosts
			allHost = append(allHost, etcdHost...)
		case "api-server":
			apiServerHost = v.Hosts
			allHost = append(allHost, apiServerHost...)
		case "scheduler":
			schedulerHost = v.Hosts
			allHost = append(allHost, schedulerHost...)
		case "controller-manager":
			controllerManagerHost = v.Hosts
			allHost = append(allHost, controllerManagerHost...)
		}
	}
	tmp := make(map[string]string)
	for _, v := range allHost {
		tmp[v] = v
	}
	for k, _ := range tmp {
		master = append(master, k)
	}

	hostInfo := info{
		Master:                master,
		EtcdHost:              etcdHost,
		ApiServerHost:         apiServerHost,
		SchedulerHost:         schedulerHost,
		ControllerManagerHost: controllerManagerHost,
		NodeHost:              nodeHost,
		HaproxyHost:           haproxyHost,
		KeepalivedHost:        keepalivedHost,
	}

	box := packr.NewBox("../template")
	hosts, _ := box.FindString("ansibleHosts/hosts")
	path := utils.Render(hostInfo, hosts, "hosts")
	return path
}
