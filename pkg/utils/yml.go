package utils

import (
	yml "gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
)

type Info struct {
	Name  string
	Hosts []string
	Dir   string
}

type Master struct {
	components []Info
}

type NodeInfo struct {
	Name string
	Dir  string
}

type Node struct {
	Hosts      []string
	Components []NodeInfo
}

type Plugin struct {
	CoreDns string
	Calico  string
}

type CIDR struct {
	ServiceCIDR string
}

type Docker struct {
	YumRepo         string
	DataRoot        string
	RegistryMirrors string
}

type Haproxy struct {
	Hosts        []string
	FrontendPort int
}

type Keepalived struct {
	Hosts     []string
	Vip       string
	Interface string
}

type Packages struct {
	DownloadDir string
	Url         []string
}

type Yaml struct {
	Certificate string
	K8s         Master
	Docker      Docker
	Plugin      Plugin
	CIDR        CIDR
	Haproxy     Haproxy
	Keepalived  Keepalived
	Packages    Packages
}

//  parser yaml to go struct
func ParserYml(filePath string) map[string]Yaml {
	file, fileErr := ioutil.ReadFile(filePath)
	if fileErr != nil {
		log.Fatalln("open %s error", filePath)
	}

	data := make(map[string]Yaml)

	unmarshalErr := yml.Unmarshal(file, &data)
	if unmarshalErr != nil {
		log.Fatalln("can not  unmarshal %s", filePath)
	}

	return data
}
