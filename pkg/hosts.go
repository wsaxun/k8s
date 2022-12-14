package pkg

import (
	"embed"
	"k8s/pkg/utils"
	"path/filepath"
	"regexp"
)

func Inventory(config utils.Config, fs embed.FS) string {
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

	hosts, _ := fs.ReadFile("template/ansibleHosts/hosts")
	path := utils.Render(hostInfo, string(hosts), "hosts")
	return path
}

func IncrementInventory(config utils.Config, fs embed.FS) string {
	type info struct {
		IncrementHost []string
	}

	nodeHost := getExistNode(config.Packages.DownloadDir)
	incrementHost := diffArray(config.K8s.Node.Hosts, nodeHost)
	hostInfo := info{
		IncrementHost: incrementHost,
	}

	hosts, _ := fs.ReadFile("template/ansibleHosts/incrementHosts")
	path := utils.Render(hostInfo, string(hosts), "incrementHosts")
	return path
}

func unrepeatedArray(slc []string) []string {
	result := []string{}
	tempMap := map[string]string{} // ?????????????????????
	for _, e := range slc {
		l := len(tempMap)
		tempMap[e] = e
		if len(tempMap) != l { // ??????map??????map?????????????????????????????????
			result = append(result, e)
		}
	}
	return result
}

func getExistNode(downloadDir string) []string {
	cmd := filepath.Join(downloadDir, "kubectl") + ` get node -o wide|grep -oE "[ ]{1,}([0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3})[ ]{1,}"`
	cmdResult := utils.Cmd("bash", "-c", cmd)
	r, _ := regexp.Compile(`[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}`)
	ipArray := r.FindAllString(cmdResult, -1)
	return ipArray
}

func diffArray(a []string, b []string) []string {
	var diffArray []string
	temp := map[string]struct{}{}

	for _, val := range b {
		if _, ok := temp[val]; !ok {
			temp[val] = struct{}{}
		}
	}

	for _, val := range a {
		if _, ok := temp[val]; !ok {
			diffArray = append(diffArray, val)
		}
	}

	return diffArray
}
