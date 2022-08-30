package pkg

import (
	"github.com/gobuffalo/packr"
	"k8s/pkg/utils"
	"os"
	"path/filepath"
)

func ConfigCsr(cache string, duration string) {
	pkiPath := filepath.Join(cache, "/kubernetes/pki/")
	_, err := os.Stat(pkiPath)
	if err != nil {
		_ := os.MkdirAll(pkiPath, 0755)
	}

	var Info struct {
		Duration string
	}
	Info.Duration = duration
	box := packr.NewBox("../template/pki")

	var files = []string{"admin-csr.json", "apiserver-csr.json", "ca-config.json", "ca-csr.json", "etcd-ca-csr.json",
		"etcd-csr.json", "front-proxy-ca-csr.json", "front-proxy-client-csr.json", "kube-proxy-csr.json",
		"manager-csr.json", "scheduler-csr.json"}

	for _, v := range files {
		context, _ := box.FindString(v)
		path := utils.Render(Info, context, v)
		_ := os.Rename(path, filepath.Join(pkiPath, v))

	}
}
