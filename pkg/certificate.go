package pkg

import (
	"github.com/gobuffalo/packr"
	"k8s/pkg/utils"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

func ConfigCsr(cache string, duration string) {
	pkiPath := filepath.Join(cache, "/kubernetes/csr/")
	_, err := os.Stat(pkiPath)
	if err != nil {
		err = os.MkdirAll(pkiPath, 0755)
		if err != nil {
			log.Fatalln(err)
		}
	}

	var Info struct {
		Duration string
	}
	Info.Duration = duration
	box := packr.NewBox("../template/csr")

	var files = []string{"admin-csr.json", "apiserver-csr.json", "ca-config.json", "ca-csr.json", "etcd-ca-csr.json",
		"etcd-csr.json", "front-proxy-ca-csr.json", "front-proxy-client-csr.json", "kube-proxy-csr.json",
		"manager-csr.json", "scheduler-csr.json"}

	for _, v := range files {
		context, _ := box.FindString(v)
		path := utils.Render(Info, context, v)
		err := os.Rename(path, filepath.Join(pkiPath, v))
		if err != nil {
			log.Fatalln(err)
		}

	}
}

func Cert(downloadCache string, etcdhost, allHost string, inventory string) {
	fileName := "generateCert.yml"
	box := packr.NewBox("../template")
	certYml, _ := box.FindString(fileName)
	type info struct {
		DownloadDir string
		Allhost     string
		Etcdhost    string
	}
	content := info{DownloadDir: downloadCache, Allhost: allHost, Etcdhost: etcdhost}
	path := utils.Render(content, certYml, fileName)
	utils.Playbook(path, inventory)
}

func ApiServerCertHost(config utils.Config) string {
	var hosts string = "127.0.0.1,kubernetes,kubernetes.default,kubernetes.default.svc,kubernetes.default.svc.cluster,kubernetes.default.svc.cluster.local"
	vip := config.Keepalived.Vip
	hosts = hosts + "," + vip
	for _, v := range config.K8s.Master.Components {
		if v.Name == "etcd" {
			continue
		}
		for _, host := range v.Hosts {
			if !strings.Contains(hosts, host) {
				hosts = hosts + "," + host
			}
		}
	}
	cidr := config.K8s.CIDR.ServiceCIDR
	ip := startIP(cidr)
	hosts = hosts + "," + ip
	return hosts
}

func EtcdHost(config utils.Config) string {
	var hosts string = "127.0.0.1"
	for _, v := range config.K8s.Master.Components {
		if v.Name != "etcd" {
			continue
		}
		for _, host := range v.Hosts {
			hosts = hosts + "," + host
		}
	}
	return hosts
}

func startIP(cidr string) string {
	ipReg := `[0-9]{1,}\.[0-9]{1,}\.[0-9]{1,}\.[0-9]{1,}`
	r, _ := regexp.Compile(ipReg)
	ipStr := r.FindString(cidr)
	tmp := strings.Split(ipStr, ".")
	a, _ := strconv.Atoi(tmp[3])
	a += 1
	ip := tmp[0] + "." + tmp[1] + "." + tmp[2] + "." + strconv.Itoa(a)
	return ip
}
