package pkg

import (
	"github.com/gobuffalo/packr"
	"k8s/pkg/utils"
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
	path := utils.Render(urlInfo, downloadYml, fileName)
	utils.Playbook(path, inventory)
}
