package biz

import (
	"os"

	"gopkg.in/yaml.v3"
)

type config struct {
	YAPI yapiConf `yaml:"yapi"`
	MS   msConf   `yaml:"ms"`
}

type yapiConf struct {
	Host    string `yaml:"host"`
	Project string `yaml:"project"`
	Token   string `yaml:"token"`
}

type msConf struct {
	Host        string `yaml:"host"`
	Workspace   string `yaml:"workspace"`
	Project     string `yaml:"project"`
	Application string `yaml:"application"`
	AccessKey   string `yaml:"accessKey"`
	SecretKey   string `yaml:"secretKey"`
}

// read configure
func readConfig() (*config, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	apidoc := pwd + "/.apidoc.yaml"
	_, err = os.Stat(apidoc)
	if err != nil {
		return nil, err
	}
	f, err := os.Open(apidoc)
	if err != nil {
		return nil, err
	}
	conf := &config{}
	err = yaml.NewDecoder(f).Decode(conf)
	if err != nil {
		return nil, err
	}
	return conf, nil
}
