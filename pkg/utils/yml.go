package utils

import (
	yml "gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
)

// cache
var AnsibleCache = "/tmp/"

// Default
type Config struct {
	K8s        K8s        `yaml:"k8s"`
	Haproxy    Haproxy    `yaml:"haproxy"`
	Keepalived Keepalived `yaml:"keepalived"`
	Packages   Packages   `yaml:"packages"`
	CRI        CRI        `yaml:"cri"`
}

// CRI
type CRI struct {
	Containerd Containerd `yaml:"containerd"`
	Docker     Docker     `yaml:"docker"`
}

// Components
type Components struct {
	Name    string   `yaml:"name"`
	Hosts   []string `yaml:"hosts"`
	Dir     string   `yaml:"dir"`
	DataDir string   `yaml:"dataDir"`
}

// Node
type Node struct {
	Hosts      []string         `yaml:"hosts"`
	Components []Nodecomponents `yaml:"components"`
}

// Plugin
type Plugin struct {
	Name       string `yaml:"name"`
	Dns        string `yaml:"dns"`
	PodCIDR    string `yaml:"podCIDR"`
	CalicoUrl  string `yaml:"calicoUrl"`
	FlannelUrl string `yaml:"flannelUrl"`
}

// CIDR
type CIDR struct {
	ServiceCIDR string `yaml:"serviceCIDR"`
}

// Keepalived
type Keepalived struct {
	Hosts     []string `yaml:"hosts"`
	Vip       string   `yaml:"vip"`
	Interface string   `yaml:"interface"`
}

// Packages
type Packages struct {
	DownloadDir string   `yaml:"downloadDir"`
	Url         []string `yaml:"url"`
}

// Master
type Master struct {
	Components []Components `yaml:"components"`
}

// Nodecomponents
type Nodecomponents struct {
	Name string `yaml:"name"`
	Dir  string `yaml:"dir"`
}

// Docker
type Docker struct {
	YumRepo          string `yaml:"yumRepo"`
	DataRoot         string `yaml:"dataRoot"`
	RegistryMirrors  string `yaml:"registryMirrors"`
	PodInfraCtrImage string `yaml:"pod-infra-container-image"`
}

// Containerd
type Containerd struct {
	YumRepo         string `yaml:"yumRepo"`
	DataRoot        string `yaml:"dataRoot"`
	RegistryMirrors string `yaml:"registryMirrors"`
	SandboxImage    string `yaml:"sandboxImage"`
}

// Haproxy
type Haproxy struct {
	Hosts        []string `yaml:"hosts"`
	FrontendPort int      `yaml:"frontendPort"`
}

// K8s
type K8s struct {
	Master      Master   `yaml:"master"`
	Node        Node     `yaml:"node"`
	Plugin      []Plugin `yaml:"plugin"`
	CIDR        CIDR     `yaml:"CIDR"`
	Certificate string   `yaml:"certificate"`
}

//  parser yaml to go struct
func ParserYml(filePath string) Config {
	file, fileErr := ioutil.ReadFile(filePath)
	if fileErr != nil {
		log.Fatalln("open %s error", filePath)
	}

	var config Config

	unmarshalErr := yml.Unmarshal(file, &config)
	if unmarshalErr != nil {
		log.Fatalln("can not  unmarshal %s", filePath)
	}

	return config
}
