package main

import (
	"k8s/pkg"
	"k8s/pkg/utils"
	"strconv"
)

func main() {
	// parser install yml
	config := utils.ParserYml("./configs/install.yml")
	cache := config.Packages.DownloadDir
	urls := config.Packages.Url

	// download
	utils.Download(cache, urls)

	// generate cert
	pkg.ConfigCsr(cache, config.K8s.Certificate)
	allHost := pkg.ApiServerCertHost(config)
	etcdHost := pkg.EtcdHost(config)
	pkg.Cert(cache, etcdHost, allHost)

	// generate kubeconfig
	pkg.KubeConfig(cache, config.Keepalived.Vip, config.Haproxy.FrontendPort)

	// install haproxy
	haproxyHost := make(map[string]string)
	for k, v := range config.Haproxy.Hosts {
		haproxyHost["haproxy"+strconv.Itoa(k+1)] = v
	}

	haproxy := pkg.Haproxy{
		Port:     config.Haproxy.FrontendPort,
		HostInfo: haproxyHost,
	}

	haproxy.InstallHaproxy("127.0.0.1")
	// init env
	//pkg.InitMasterEnv("127.0.0.1")
	//pkg.InitNodeEnv("127.0.0.1")
	//
	//docker := pkg.Docker{
	//	YumRepo:         "http://mirrors.aliyun.com/docker-ce/linux/centos/docker-ce.repo",
	//	DataRoot:        "/var/lib/docker",
	//	RegistryMirrors: "https://mvaav0ar.mirror.aliyuncs.com",
	//}
	//docker.InstallDocker("127.0.0.1")
}
