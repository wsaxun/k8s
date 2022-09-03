package pkg

import (
	"github.com/gobuffalo/packr"
	"k8s/pkg/utils"
)

func Inventory(config utils.Config) string {
	type info struct {
		Kubernetes            []string
		Master                []string
		EtcdHost              []string
		ApiServerHost         []string
		SchedulerHost         []string
		ControllerManagerHost []string
		NodeHost              []string
		HaproxyHost           []string
		KeepalivedHost        []string
	}

	var kubernetes []string
	var tmpMaster []string
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
			tmpMaster = append(tmpMaster, etcdHost...)
		case "api-server":
			apiServerHost = v.Hosts
			tmpMaster = append(tmpMaster, apiServerHost...)
		case "scheduler":
			schedulerHost = v.Hosts
			tmpMaster = append(tmpMaster, schedulerHost...)
		case "controller-manager":
			controllerManagerHost = v.Hosts
			tmpMaster = append(tmpMaster, controllerManagerHost...)
		}
	}
	// master
	master = unrepeatedArray(tmpMaster)
	// kubernets = master + node
	tmpK8s := append(master, nodeHost...)
	kubernetes = unrepeatedArray(tmpK8s)

	hostInfo := info{
		Kubernetes:            kubernetes,
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

func unrepeatedArray(slc []string) []string {
	result := []string{}
	tempMap := map[string]string{} // 存放不重复主键
	for _, e := range slc {
		l := len(tempMap)
		tempMap[e] = e
		if len(tempMap) != l { // 加入map后，map长度变化，则元素不重复
			result = append(result, e)
		}
	}
	return result
}
