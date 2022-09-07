package pkg

import (
	"embed"
	"k8s/pkg/utils"
)

type Bootstrap struct {
	utils.K8SSoftware
	DownloadDir string
	Vip         string
	Port        int
}

func (b *Bootstrap) Install(host string, inventory string, fs embed.FS) {
	tokenId, token := utils.Token()
	b.config(tokenId, token, fs)
	ymlName := "bootstrap.yml"
	yml, _ := fs.ReadFile("template/bootstrap.yml")

	type info struct {
		DownloadDir string
		Vip         string
		Port        int
		TokenId     string
		Token       string
	}
	bootstrapInfo := info{
		DownloadDir: b.DownloadDir,
		Vip:         b.Vip,
		Port:        b.Port,
		TokenId:     tokenId,
		Token:       token,
	}
	path := utils.Render(bootstrapInfo, string(yml), ymlName)
	utils.Playbook(path, inventory)
}

func (b *Bootstrap) config(tokenId, token string, fs embed.FS) {
	context, _ := fs.ReadFile("template/softwareConfig/bootstrap.secret.yml")

	type info struct {
		TokenId string
		Token   string
	}
	bootstrapInfo := info{
		TokenId: tokenId,
		Token:   token,
	}
	utils.Render(bootstrapInfo, string(context), "bootstrap.secret.yml")
}
