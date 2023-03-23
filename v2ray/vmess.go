package v2ray

import (
	"Surfing/util"

	"encoding/base64"
	"encoding/json"
	"strings"
	"time"
)

type Config struct {
	Add  string `json:"add"`
	Aid  string `json:"aid"`
	Alpn string `json:"alpn"`
	Fp   string `json:"fp"`
	Host string `json:"host"`
	Id   string `json:"id"`
	Net  string `json:"net"`
	Path string `json:"path"`
	Port string `json:"port"`
	Ps   string `json:"ps"`
	Scy  string `json:"scy"`
	Sni  string `json:"sni"`
	Tls  string `json:"tls"`
	Type string `json:"type"`
	V    string `json:"v"`
}

func VmessToSubscription(vmessUrls []string) string {

	var subscription string

	for _, url := range vmessUrls {
		subscription += url + "\n" // add a newline character after each URL
	}

	subscription = base64.StdEncoding.EncodeToString([]byte(subscription))

	util.PrintLog("Update subscription :" + subscription)

	return subscription
}

func CreateVmessFromVmess(cloudfrontDomainName string, originalVmess string) string {

	vmessJson := VmessURLDecode(originalVmess)
	config := ParseVmessJson(vmessJson)
	config = ChangeDomainName(cloudfrontDomainName, config)
	newVmessJson := ConfigObjectToString(config)
	newVmessURL := EncodingToVmesssURL(newVmessJson)

	util.PrintLog("Created new Vmess URL using :" + cloudfrontDomainName)
	return newVmessURL

}

func VmessURLDecode(vmessString string) string {
	decodedString, err := base64.URLEncoding.DecodeString(strings.TrimPrefix(vmessString, "vmess://"))

	if err != nil {
		util.PrintLog("Error decoding base64 string:" + err.Error())
	}

	return string(decodedString)
}

func ParseVmessJson(jsonString string) Config {
	var config Config
	err := json.Unmarshal([]byte(jsonString), &config)
	if err != nil {
		util.PrintLog("Error parsing JSON:" + err.Error())
	}

	return config
}

func ChangeDomainName(domainName string, config Config) Config {
	config.Add = domainName
	config.Host = domainName
	config.Sni = domainName
	config.Ps = time.Now().Format("2006-01-02 15:04:05")

	return config
}

func ConfigObjectToString(config Config) string {
	jsonString, err := json.Marshal(config)
	if err != nil {
		util.PrintLog("Error encoding JSON:" + err.Error())
	}

	return string(jsonString)
}

func EncodingToVmesssURL(jsonString string) string {
	base64String := base64.URLEncoding.EncodeToString([]byte(jsonString))

	return "vmess://" + base64String
}
