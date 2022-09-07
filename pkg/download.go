package pkg

import (
	"embed"
	"k8s/pkg/utils"
)

type Software struct {
	DownloadPackage string
	URL             []string
}

func (s *Software) DownloadPackages(inventory string, fs embed.FS) {
	type info struct {
		DownloadPackage string
		URL             []string
	}

	urlInfo := info{DownloadPackage: s.DownloadPackage, URL: s.URL}

	fileName := "download.yml"
	downloadYml, _ := fs.ReadFile("template/download.yml")
	path := utils.Render(urlInfo, string(downloadYml), fileName)
	utils.Playbook(path, inventory)
}
