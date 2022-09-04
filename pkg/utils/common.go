package utils

import (
	"crypto/md5"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"
)

func GetCache() string {
	return AnsibleCache
}

func Render(data interface{}, templateStr, fileName string) (path string) {
	tmp := template.New(fileName)
	tmp.Funcs(template.FuncMap{"GetCache": GetCache})
	tmpl, err := tmp.Parse(templateStr)
	if err != nil {
		log.Fatalln(err)
	}
	file, err := os.OpenFile(filepath.Join(AnsibleCache, fileName), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()
	err = tmpl.Execute(file, data)
	if err != nil {
		log.Fatalln(err)
	}
	return filepath.Join(AnsibleCache, fileName)
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
