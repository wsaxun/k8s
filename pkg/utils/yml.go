package utils

import (
	yml "gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
)

//  parser yaml to go struct
func ParserYml(filePath string) map[interface{}]interface{} {
	file, fileErr := ioutil.ReadFile(filePath)
	if fileErr != nil {
		log.Fatalln("open %s error", filePath)
	}

	data := make(map[interface{}]interface{})

	unmarshalErr := yml.Unmarshal(file, &data)
	if unmarshalErr != nil {
		log.Fatalln("can not  unmarshal %s", filePath)
	}

	return data
}
