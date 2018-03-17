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

	DbHost     string `json:"DbHost"`
	DbPort     int    `json:"DbPort"`
	DbUsername string `json:"DbUsername"`
	DbPassword string `json:"DbPassword"`
	DbName     string `json:"DbName"`

	DirUpload string `json:"DirUpload"`
	DirStatic string `json:"DirStatic"`

	QQConnectAppId    string `json:"QQConnectAppId"`
	QQConnectSecret   string `json:"QQConnectSecret"`
	QQConnectRedirect string `json:"QQConnectRedirect"`
}

const (
	configFile = "./conf/config.json"
)

func init() {
	Conf = readConfig()
	log.Println("success read config")
}

func readConfig() *Config {
	file, err := os.Open(configFile)
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
