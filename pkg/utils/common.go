package utils

import (
	"crypto/md5"
	"fmt"
	"github.com/jessevdk/go-flags"
	uuid "github.com/satori/go.uuid"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
)

type K8SSoftware struct{}

func (s *K8SSoftware) Install() {
	log.Fatalln("you must implement this function")
}

func (s *K8SSoftware) config() {
	log.Fatalln("you must implement this function")
}

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
	defer func() {
		if err := file.Close(); err != nil {
			log.Fatalln(err)
		}
	}()
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
		log.Println(name, args, string(out))
		log.Fatalln(err)
	}
	return string(out)
}

func FileName(url string) string {
	tmp := strings.Split(url, "/")
	length := len(tmp)
	name := tmp[length-1]
	return name
}

type CmdOption struct {
	PrintDefault bool   `short:"p" long:"PrintDefault" description:"print install default config"`
	InstallType  string `short:"i" long:"install" description:"k8s or node"`
	ConfigFile   string `short:"f" long:"file" description:"install config file"`
}

func CmdArgs() CmdOption {
	var opt CmdOption
	flags.Parse(&opt)
	return opt
}

func PathIsExist(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	return false
}
