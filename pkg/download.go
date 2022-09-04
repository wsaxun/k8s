package pkg

import (
	"github.com/gobuffalo/packr"
	"k8s/pkg/utils"
	"os"
	"strings"
)

type Software struct {
	DownloadPackage string
	URL             []string
}

func (s *Software) DownloadPackages(inventory string) {
	type info struct {
		DownloadPackage string
		URL             []string
	}

	urlInfo := info{DownloadPackage: s.DownloadPackage, URL: s.URL}

	fileName := "download.yml"
	box := packr.NewBox("../template")
	downloadYml, _ := box.FindString(fileName)
	utils.Render(urlInfo, downloadYml, fileName)
	//path := utils.Render(urlInfo, downloadYml, fileName)
	//utils.Playbook(path, inventory)
}

func (s *Software) IsDownload(url string) bool {
	tmp := strings.Split(url, "/")
	length := len(tmp)
	softwareName := tmp[length-1]
	_, err := os.Stat(s.DownloadPackage + "/" + softwareName)
	if err == nil {
		return false
	}
	return true
}
