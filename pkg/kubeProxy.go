package pkg

import (
	"embed"
	"k8s/pkg/utils"
)

type Proxy struct {
	utils.K8SSoftware
	Vip         string
	Port        int
	Dir         string
	DownloadDir string
	PodCIDR     string
}

func (p *Proxy) Install(host string, inventory string, fs embed.FS) {
	p.config(inventory, fs)
	ymlName := "kubeProxy.yml"
	yml, _ := fs.ReadFile("template/" + ymlName)
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
	path := utils.Render(kubeProxyInfo, string(yml), ymlName)
	utils.Playbook(path, inventory)
}

func (p *Proxy) config(inventory string, fs embed.FS) {
	//p.kubeConfig(inventory, fs, genKubConfFlag)
	context, _ := fs.ReadFile("template/softwareConfig/kube-proxy.service")
	utils.Render(p, string(context), "kube-proxy.service")
	service, _ := fs.ReadFile("template/softwareConfig/kube-proxy.config.yml")
	utils.Render(p, string(service), "kube-proxy.config.yml")
}

//func (p *Proxy) kubeConfig(inventory string, fs embed.FS, genKubConfFlag bool) {
//	if !genKubConfFlag {
//		return
//	}
//	cmd := p.DownloadDir + "/kubectl -n kube-system create serviceaccount kube-proxy"
//	utils.Cmd("bash", "-c", cmd)
//	time.Sleep(30 * time.Second)
//	cmd = p.DownloadDir + "/kubectl create clusterrolebinding system:kube-proxy --clusterrole system:node-proxier --serviceaccount kube-system:kube-proxy"
//	utils.Cmd("bash", "-c", cmd)
//	time.Sleep(90 * time.Second)
//	cmd = p.DownloadDir + `/kubectl -n kube-system get sa/kube-proxy --output=jsonpath='{.secrets[0].name}'`
//	secrete := utils.Cmd("bash", "-c", cmd)
//	cmd = p.DownloadDir + `/kubectl -n kube-system get secret/` + secrete + `   --output=jsonpath='{.data.token}' | /usr/bin/base64 -d`
//	token := utils.Cmd("bash", "-c", cmd)
//	type info struct {
//		DownloadDir string
//		Vip         string
//		Port        int
//		Token       string
//	}
//	proxyInfo := info{
//		DownloadDir: p.DownloadDir,
//		Vip:         p.Vip,
//		Port:        p.Port,
//		Token:       token,
//	}
//	context, _ := fs.ReadFile("template/kubeProxyKubeConfig.yml")
//	path := utils.Render(proxyInfo, string(context), "kubeProxyKubeConfig.yml")
//	utils.Playbook(path, inventory)
//}
