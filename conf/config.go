package conf

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

var Conf *Config

type Config struct {
	SiteName    string `json:"SiteName"`
	SiteAddr    string `json:"SiteAddr"`
	SitePort    string `json:"SitePort"`
	SitePortSSL string `json:"SitePortSSL"`
	SecretKey   string `json:"SecretKey"`
	DbUsername  string `json:"DbUsername"`
	DbPassword  string `json:"DbPassword"`
	DbName      string `json:"DbName"`
}

func init() {
	Conf = readConfig()
	log.Println("success read config")
}

func readConfig() *Config {
	if Conf != nil {
		return Conf
	}
	inputFile := "./conf/config.json"
	buf, err := ioutil.ReadFile(inputFile)
	if err != nil {
		panic(err.Error())
	}

	log.Println(string(buf))

	var c *Config = &Config{}
	err = json.Unmarshal(buf, c)

	if err != nil {
		panic(err.Error())
	}

	return c
}
