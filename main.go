package main

import (
	"k8s/pkg"
	"k8s/pkg/utils"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {
	// parser install yml
	config := utils.ParserYml("./configs/install.yml")
	softwareDownloadDir := config.Packages.DownloadDir
	urls := config.Packages.Url
	utils.AnsibleCache = filepath.Join(softwareDownloadDir, "ansibleCache")
	err := os.MkdirAll(utils.AnsibleCache, 0755)
	if err != nil {
		log.Fatalln(err)
	}

	var etcdName string
	for _, v := range urls {
		if strings.Index(v, "etcd") >= 1 {
			etcdName = utils.FileName(v)
		}
	}

	// generate ansible inventory
	inventory := pkg.Inventory(config)

	// download
	// TODO
	software := pkg.Software{
		DownloadPackage: softwareDownloadDir,
		URL:             urls,
	}
	software.DownloadPackages(inventory)

	// generate cert
	pkg.ConfigCsr(softwareDownloadDir, config.K8s.Certificate)
	allHost := pkg.ApiServerCertHost(config)
	etcdHost := pkg.EtcdHost(config)
	pkg.Cert(softwareDownloadDir, etcdHost, allHost, inventory)

	// generate kubeconfig
	pkg.KubeConfig(softwareDownloadDir, config.Keepalived.Vip, config.Haproxy.FrontendPort, inventory)

	// init env
	// TODO
	pkg.InitMasterEnv("master", inventory)
	pkg.InitNodeEnv("node", inventory)

	// install docker
	yumRepo := config.Docker.YumRepo
	dataRoot := config.Docker.DataRoot
	mirror := config.Docker.RegistryMirrors
	docker := pkg.Docker{
		YumRepo:         yumRepo,
		DataRoot:        dataRoot,
		RegistryMirrors: mirror,
	}
	docker.Install("master", inventory)
	docker.Install("node", inventory)

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
	haproxy.Install("haproxy", inventory)

	// install keepalived
	// TODO
	keepalived := pkg.Keepalived{
		Interface: config.Keepalived.Interface,
		Host:      config.Keepalived.Hosts,
		Vip:       config.Keepalived.Vip,
	}
	keepalived.Install("keepalived", inventory)

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
		EtcdName:    etcdName,
	}
	etcd.Install("etcd", inventory)

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
		DownloadDir: softwareDownloadDir,
		ServiceCIDR: config.K8s.CIDR.ServiceCIDR,
		EtcdHost:    etcdHostArray,
	}
	apiserver.Install("api-server", inventory)

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
	contr.Install("controller-manager", inventory)

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
	scheduler.Install("scheduler", inventory)

	// install bootstrap
	// TODO
	bootstrap := pkg.Bootstrap{
		DownloadDir: config.Packages.DownloadDir,
		Vip:         config.Keepalived.Vip,
		Port:        config.Haproxy.FrontendPort,
	}
	bootstrap.Install("127.0.0.1", inventory)

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
		DownloadDir: softwareDownloadDir,
	}
	kubelet.Install("kubernetes", inventory)

	// install kube-proxy
	// TODO
	var kubeProxyDir string
	for _, v := range config.K8s.Node.Components {
		if v.Name == "kubproxy" {
			kubeProxyDir = v.Dir
		}
	}
	kubeProxy := pkg.Proxy{
		Vip:         config.Keepalived.Vip,
		Port:        config.Haproxy.FrontendPort,
		Dir:         kubeProxyDir,
		DownloadDir: softwareDownloadDir,
		PodCIDR:     podCIDR,
	}
	kubeProxy.Install("kubernetes", inventory)

	// install calico
	// TODO
	var calicoUrl string
	for _, v := range config.K8s.Plugin {
		if v.Name == "calico" {
			calicoUrl = v.CalicoUrl
		}
	}
	calico := pkg.Calico{
		DownloadDir: softwareDownloadDir,
		Url:         calicoUrl,
		PodCIDR:     podCIDR,
	}
	calico.Install()

	// install coredns
	// TODO
	coredns := pkg.CoreDns{
		Dns:         dns,
		DownloadDir: softwareDownloadDir,
	}
	coredns.Install()
}
