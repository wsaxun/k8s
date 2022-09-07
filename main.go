package main

import (
	"embed"
	"fmt"
	"github.com/jessevdk/go-flags"
	"k8s/pkg"
	"k8s/pkg/utils"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

//go:embed configs
var CONFIGS embed.FS

//go:embed template
var TEMPLATE embed.FS

func main() {
	// parser cmd option
	var cmdOption utils.CmdOption
	flags.Parse(&cmdOption)

	if cmdOption.PrintDefault {
		context, _ := CONFIGS.ReadFile("configs/install.yml")
		fmt.Println(string(context))
		os.Exit(0)
	}

	if !utils.PathIsExist(cmdOption.ConfigFile) {
		log.Fatalln("file is not exist")
	}
	log.Println("k8s start install")

	// parser install yml
	config := utils.ParserYml(cmdOption.ConfigFile)
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

	//
	yumRepo := config.Docker.YumRepo
	dataRoot := config.Docker.DataRoot
	mirror := config.Docker.RegistryMirrors
	docker := pkg.Docker{
		YumRepo:         yumRepo,
		DataRoot:        dataRoot,
		RegistryMirrors: mirror,
	}
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
	var kubeProxyDir string
	for _, v := range config.K8s.Node.Components {
		if v.Name == "kubproxy" {
			kubeProxyDir = v.Dir
		}
	}
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

	//generate ansible inventory
	log.Println("generate ansible hosts")
	inventory := pkg.Inventory(config, TEMPLATE)

	// node 部署
	if cmdOption.InstallType == "node" {
		inventory := pkg.IncrementInventory(config, TEMPLATE)

		log.Println("init node computer")
		pkg.InitNodeEnv("increment", inventory, TEMPLATE)

		log.Println("install docker")
		docker.Install("increment", inventory, TEMPLATE)

		log.Println("install kubelet")
		kubelet := pkg.Kubelet{
			CoreDns:     dns,
			Dir:         kubeletDir,
			DownloadDir: softwareDownloadDir,
		}
		kubelet.Install("increment", inventory, TEMPLATE)

		// install kube-proxy
		log.Println("install kube-haproxy")
		kubeProxy := pkg.Proxy{
			Vip:         config.Keepalived.Vip,
			Port:        config.Haproxy.FrontendPort,
			Dir:         kubeProxyDir,
			DownloadDir: softwareDownloadDir,
			PodCIDR:     podCIDR,
		}
		kubeProxy.Install("increment", inventory, TEMPLATE)
		os.Exit(0)
	}

	// download
	log.Println("download software")

	software := pkg.Software{
		DownloadPackage: softwareDownloadDir,
		URL:             urls,
	}
	software.DownloadPackages(inventory, TEMPLATE)

	// generate cert
	log.Println("generate cert")

	pkg.ConfigCsr(softwareDownloadDir, config.K8s.Certificate, TEMPLATE)
	allHost := pkg.ApiServerCertHost(config)
	etcdHost := pkg.EtcdHost(config)
	pkg.Cert(softwareDownloadDir, etcdHost, allHost, inventory, TEMPLATE)

	// generate kubeconfig
	log.Println("generate kubconfig")
	pkg.KubeConfig(softwareDownloadDir, config.Keepalived.Vip, config.Haproxy.FrontendPort, inventory, TEMPLATE)

	// init env
	log.Println("init master computer")
	pkg.InitMasterEnv("master", inventory, TEMPLATE)
	log.Println("init node computer")
	pkg.InitNodeEnv("node", inventory, TEMPLATE)

	// install docker
	log.Println("install docker")
	docker.Install("master", inventory, TEMPLATE)
	docker.Install("node", inventory, TEMPLATE)

	// install haproxy
	log.Println("install haproxy")

	haproxyHost := make(map[string]string)
	for k, v := range config.Haproxy.Hosts {
		haproxyHost["haproxy"+strconv.Itoa(k+1)] = v
	}
	haproxy := pkg.Haproxy{
		Port:     config.Haproxy.FrontendPort,
		HostInfo: haproxyHost,
	}
	haproxy.Install("haproxy", inventory, TEMPLATE)

	// install keepalived
	log.Println("install keepalived")

	keepalived := pkg.Keepalived{
		Interface: config.Keepalived.Interface,
		Host:      config.Keepalived.Hosts,
		Vip:       config.Keepalived.Vip,
	}
	keepalived.Install("keepalived", inventory, TEMPLATE)

	// install etcd
	log.Println("install etcd")

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
	etcd.Install("etcd", inventory, TEMPLATE)

	// install apiServer
	log.Println("install api-server")

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
	apiserver.Install("api-server", inventory, TEMPLATE)

	// install controllerManager
	log.Println("install controllerManager")

	contr := pkg.ControllerManager{
		Dir:         contrDir,
		PodCIDR:     podCIDR,
		DownloadDir: config.Packages.DownloadDir,
	}
	contr.Install("controller-manager", inventory, TEMPLATE)

	// install scheduler
	log.Println("install scheduler")

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
	scheduler.Install("scheduler", inventory, TEMPLATE)

	// install bootstrap
	log.Println("config bootstrap")

	bootstrap := pkg.Bootstrap{
		DownloadDir: config.Packages.DownloadDir,
		Vip:         config.Keepalived.Vip,
		Port:        config.Haproxy.FrontendPort,
	}
	bootstrap.Install("127.0.0.1", inventory, TEMPLATE)

	// install kubelet
	log.Println("install kubelet")
	kubelet := pkg.Kubelet{
		CoreDns:     dns,
		Dir:         kubeletDir,
		DownloadDir: softwareDownloadDir,
	}
	kubelet.Install("kubernetes", inventory, TEMPLATE)

	// install kube-proxy
	log.Println("install kube-haproxy")
	kubeProxy := pkg.Proxy{
		Vip:         config.Keepalived.Vip,
		Port:        config.Haproxy.FrontendPort,
		Dir:         kubeProxyDir,
		DownloadDir: softwareDownloadDir,
		PodCIDR:     podCIDR,
	}
	kubeProxy.Install("kubernetes", inventory, TEMPLATE)

	// install calico
	log.Println("install calico")

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
	log.Println("install coreDns")

	coredns := pkg.CoreDns{
		Dns:         dns,
		DownloadDir: softwareDownloadDir,
	}
	coredns.Install(TEMPLATE)
}
