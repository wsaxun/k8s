package utils

import (
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
