package pkg

import (
	"github.com/gobuffalo/packr"
	"k8s/pkg/utils"
)

type Bootstrap struct {
	DownloadDir string
	Vip         string
	Port        int
}

func (b *Bootstrap) InstallBootstrap(host string, inventory string) {
	tokenId, token := utils.Token()
	b.config(tokenId, token)
	ymlName := "bootstrap.yml"
	box := packr.NewBox("../template")
	yml, _ := box.FindString(ymlName)

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
	utils.Render(bootstrapInfo, yml, ymlName)
	//path := utils.Render(content, yml, ymlName)
	//utils.Playbook(path, inventory)
}

func (b *Bootstrap) config(tokenId, token string) {
	box := packr.NewBox("../template")
	context, _ := box.FindString("softwareConfig/bootstrap.secret.yml")
	type info struct {
		TokenId string
		Token   string
	}
	bootstrapInfo := info{
		TokenId: tokenId,
		Token:   token,
	}
	utils.Render(bootstrapInfo, context, "bootstrap.secret.yml")
}
