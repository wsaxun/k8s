package pkg

import (
	"github.com/gobuffalo/packr"
	"k8s/pkg/utils"
	"path/filepath"
	"time"
)

type Proxy struct {
	utils.K8SSoftware
	Vip         string
	Port        int
	Dir         string
	DownloadDir string
	PodCIDR     string
}

func (p *Proxy) Install(host string, inventory string) {
	p.config(inventory)
	ymlName := "kubeProxy.yml"
	box := packr.NewBox("../template")
	yml, _ := box.FindString(ymlName)
	type info struct {
		Dir         string
		DownloadDir string
		Host        string
	}
	kubeProxyInfo := info{
		Dir:         p.Dir,
		DownloadDir: p.DownloadDir,
		Host:        host,
	}
	path := utils.Render(kubeProxyInfo, yml, ymlName)
	utils.Playbook(path, inventory)
}

func (p *Proxy) config(inventory string) {
	p.kubeConfig(inventory)
	box := packr.NewBox("../template")
	context, _ := box.FindString("softwareConfig/kube-proxy.service")
	utils.Render(p, context, "kube-proxy.service")
	service, _ := box.FindString("softwareConfig/kube-proxy.config.yml")
	utils.Render(p, service, "kube-proxy.config.yml")
}

func (p *Proxy) kubeConfig(inventory string) {
	if utils.PathIsExist(filepath.Join(utils.GetCache(), "kubeProxyKubeConfig.yml")) {
		return
	}
	cmd := p.DownloadDir + "/kubectl " + "-n kube-system create serviceaccount kube-proxy"
	utils.Cmd("bash", "-c", cmd)
	cmd = p.DownloadDir + "/kubectl create clusterrolebinding system:kube-proxy --clusterrole system:node-proxier --serviceaccount kube-system:kube-proxy"
	utils.Cmd("bash", "-c", cmd)
	time.Sleep(60 * time.Second)
	cmd = p.DownloadDir + `/kubectl -n kube-system get sa/kube-proxy --output=jsonpath='{.secrets[0].name}'`
	secrete := utils.Cmd("bash", "-c", cmd)
	cmd = p.DownloadDir + `/kubectl -n kube-system get secret/` + secrete + `   --output=jsonpath='{.data.token}' | /usr/bin/base64 -d`
	token := utils.Cmd("bash", "-c", cmd)
	type info struct {
		DownloadDir string
		Vip         string
		Port        int
		Token       string
	}
	proxyInfo := info{
		DownloadDir: p.DownloadDir,
		Vip:         p.Vip,
		Port:        p.Port,
		Token:       token,
	}
	box := packr.NewBox("../template")
	context, _ := box.FindString("kubeProxyKubeConfig.yml")
	path := utils.Render(proxyInfo, context, "kubeProxyKubeConfig.yml")
	utils.Playbook(path, inventory)
}
