package utils

import (
	"crypto/md5"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
)

func GetCache() string {
	return AnsibleCache
}

func IsDownload(url, downloadDir string) bool {
	tmp := strings.Split(url, "/")
	length := len(tmp)
	name := tmp[length-1]
	_, err := os.Stat(filepath.Join(downloadDir, name))
	if err == nil {
		return false
	}
	return true
}

func Render(data interface{}, templateStr, fileName string) (path string) {
	tmp := template.New(fileName)
	tmp.Funcs(template.FuncMap{"GetCache": GetCache, "IsDownload": IsDownload})
	tmpl, err := tmp.Parse(templateStr)
	if err != nil {
		log.Fatalln(err)
	}
	file, err := os.OpenFile(filepath.Join(GetCache(), fileName), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()
	err = tmpl.Execute(file, data)
	if err != nil {
		log.Fatalln(err)
	}
	return filepath.Join(GetCache(), fileName)
}

func Token() (tokenId, token string) {
	uid := uuid.NewV4().Bytes()
	md := md5.Sum(uid)
	x := md[:]
	return fmt.Sprintf("%x", x)[0:6], fmt.Sprintf("%x", x)[0:16]
}

func Cmd(name string, args ...string) string {
	command := exec.Command(name, args...)
	out, err := command.CombinedOutput()
	if err != nil {
		log.Fatalln(err)
	}
	return string(out)
}
