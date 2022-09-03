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

	// generate ansible inventory
	inventory := pkg.Inventory(config)

	// download
	// TODO
	pkg.DownloadPackages(cache, urls, inventory)

	// generate cert
	pkg.ConfigCsr(cache, config.K8s.Certificate)
	allHost := pkg.ApiServerCertHost(config)
	etcdHost := pkg.EtcdHost(config)
	pkg.Cert(cache, etcdHost, allHost, inventory)

	// generate kubeconfig
	pkg.KubeConfig(cache, config.Keepalived.Vip, config.Haproxy.FrontendPort, inventory)

	// init env
	//pkg.InitMasterEnv("master", inventory)
	//pkg.InitNodeEnv("node", inventory)

	//docker := pkg.Docker{
	//	YumRepo:         "http://mirrors.aliyun.com/docker-ce/linux/centos/docker-ce.repo",
	//	DataRoot:        "/var/lib/docker",
	//	RegistryMirrors: "https://mvaav0ar.mirror.aliyuncs.com",
	//}
	//docker.InstallDocker("master", inventory)
	//docker.InstallDocker("node", inventory)

	// install haproxy
	// TODO
	haproxyHost := make(map[string]string)
	for k, v := range config.Haproxy.Hosts {
		haproxyHost["haproxy"+strconv.Itoa(k+1)] = v
	}
	haproxy := pkg.Haproxy{
		Port:     config.Haproxy.FrontendPort,
		HostInfo: haproxyHost,
	}
	haproxy.InstallHaproxy("haproxy", inventory)

	// install keepalived
	// TODO
	keepalived := pkg.Keepalived{
		Interface: config.Keepalived.Interface,
		Host:      config.Keepalived.Hosts,
		Vip:       config.Keepalived.Vip,
	}
	keepalived.InstallKeepalived("keepalived", inventory)

	// install etcd
	// TODO
	var dataDir string
	var dir string
	var etcdHostArray []string
	for _, v := range config.K8s.Master.Components {
		if v.Name == "etcd" {
			dataDir = v.DataDir
			dir = v.Dir
			etcdHostArray = v.Hosts
		}
	}
	etcd := pkg.Etcd{
		DataDir:     dataDir,
		Host:        etcdHostArray,
		Dir:         dir,
		DownloadDir: config.Packages.DownloadDir,
	}
	etcd.InstallEtcd("etcd", inventory)

	// install apiServer
	// TODO
	var apiServerDir string
	var apiServerHostArray []string
	for _, v := range config.K8s.Master.Components {
		if v.Name == "api-server" {
			apiServerDir = v.Dir
			apiServerHostArray = v.Hosts
		}
	}
	apiserver := pkg.ApiServer{
		Host:        apiServerHostArray,
		Dir:         apiServerDir,
		DownloadDir: cache,
		ServiceCIDR: config.K8s.CIDR.ServiceCIDR,
		EtcdHost:    etcdHostArray,
	}
	apiserver.InstallApiServer("api-server", inventory)

	// install controllerManager
	// TODO
	var contrDir string
	var podCIDR string
	for _, v := range config.K8s.Master.Components {
		if v.Name == "controller-manager" {
			contrDir = v.Dir
		}
	}
	for _, v := range config.K8s.Plugin {
		if v.Name == "calico" {
			podCIDR = v.PodCIDR
		}
	}
	contr := pkg.ControllerManager{
		Dir:         contrDir,
		PodCIDR:     podCIDR,
		DownloadDir: config.Packages.DownloadDir,
	}
	contr.InstallControllerManager("controller-manager", inventory)

	// install scheduler
	// TODO
	var schedulerDir string
	for _, v := range config.K8s.Master.Components {
		if v.Name == "scheduler" {
			schedulerDir = v.Dir
		}
	}
	scheduler := pkg.Scheduler{
		Dir:         schedulerDir,
		DownloadDir: config.Packages.DownloadDir,
	}
	scheduler.InstallScheduler("scheduler", inventory)

	// install bootstrap
	// TODO
	bootstrap := pkg.Bootstrap{
		DownloadDir: config.Packages.DownloadDir,
		Vip:         config.Keepalived.Vip,
		Port:        config.Haproxy.FrontendPort,
	}
	bootstrap.InstallBootstrap("127.0.0.1", inventory)

	// install kubelet
	// TODO
	var dns string
	var kubeletDir string
	for _, v := range config.K8s.Plugin {
		if v.Name == "coreDns" {
			dns = v.Dns
		}
	}
	for _, v := range config.K8s.Node.Components {
		if v.Name == "kubelet" {
			kubeletDir = v.Dir
		}
	}
	kubelet := pkg.Kubelet{
		CoreDns:     dns,
		Dir:         kubeletDir,
		DownloadDir: cache,
	}
	kubelet.InstallKubelet("kubernetes", inventory)
}
