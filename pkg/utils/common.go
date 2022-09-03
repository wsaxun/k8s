package utils

import (
	"crypto/md5"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"log"
	"os"
	"os/exec"
	"text/template"
)

func Render(data interface{}, templateStr, fileName string) (path string) {
	tmpl, err := template.New("t").Parse(templateStr)
	defer tmpl.Clone()

	if err != nil {
		log.Fatalln(err)
	}
	file, err := os.OpenFile("/tmp/"+fileName, os.O_CREATE|os.O_WRONLY, 0644)
	defer file.Close()
	err = tmpl.Execute(file, data)
	if err != nil {
		log.Fatalln(err)
	}
	return "/tmp/" + fileName
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
