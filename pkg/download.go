package pkg

import (
	"github.com/gobuffalo/packr"
	"k8s/pkg/utils"
)

func DownloadPackages(cache string, urls []string, inventory string) {
	type info struct {
		DownloadPackage string
		URL             []string
	}

	urlInfo := info{DownloadPackage: cache, URL: urls}

	fileName := "download.yml"
	box := packr.NewBox("../template")
	downloadYml, _ := box.FindString(fileName)
	utils.Render(urlInfo, downloadYml, fileName)
	//path := utils.Render(urlInfo, downloadYml, fileName)
	//utils.Playbook(path, inventory)
}
