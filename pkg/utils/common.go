package utils

import (
	"fmt"
	"log"
	"os"
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

func Download(downloadDir string, url []string) {
	Exec("127.0.0.1", "shell", "mkdir -p "+downloadDir)
	for _, v := range url {
		Exec("127.0.0.1", "shell", fmt.Sprintf("cd %s && wget %s ", downloadDir, v))
	}
	Exec("127.0.0.1", "shell", "cd %s && chmod +x *")
}
