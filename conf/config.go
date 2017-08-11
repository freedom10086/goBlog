package conf

import (
	"encoding/json"
	"log"
	"os"
)

var Conf *Config

type Config struct {
	SiteName      string `json:"SiteName"`
	SiteIpAddr    string `json:"SiteIpAddr"`
	SitePort      string `json:"SitePort"`
	SitePortSSL   string `json:"SitePortSSL"`
	SiteStaticDir string `json:"SiteStaticDir"`
	SecretKey     string `json:"SecretKey"`
	DbUsername    string `json:"DbUsername"`
	DbPassword    string `json:"DbPassword"`
	DbName        string `json:"DbName"`
	DirTemplate   string `json:"DirTemplate"`
	DirUpload   string `json:"DirUpload"`
	DirStatic     string `json:"DirStatic"`
}

const (
	dir_config = "./conf/config.json"
)

func init() {
	Conf = readConfig()
	log.Println("success read config")
}

func readConfig() *Config {
	file, err := os.Open(dir_config)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	conf := &Config{}
	err = json.NewDecoder(file).Decode(conf)

	if err != nil {
		log.Fatal(err)
	}
	return conf
}
