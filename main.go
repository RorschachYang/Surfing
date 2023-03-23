package main

import (
	"Surfing/cdn_cf"
	"Surfing/util"
	"Surfing/v2ray"
	"flag"
)

func main() {

	var arg string
	var configDir string
	flag.StringVar(&arg, "job", "", "what job to run")
	flag.StringVar(&configDir, "f", "", "where config is")

	flag.Parse()

	util.SetConfigDir(configDir)

	// 根据选项执行不同的方法
	switch arg {
	case "create":
		createSubscription()
	case "disable":
		disableDistribution()
	case "delete":
		deleteDistribution()
	default:
		util.PrintLog("Unknown job")
	}

}

func createSubscription() {

	V2rayDomainName := util.GetConfigString("V2ray.domainName")
	oldVmess := util.GetConfigString("V2ray.originalVmess")
	vmessNum := util.GetConfigInt("V2ray.vmessNum")
	dir := util.GetConfigString("Nginx.subscriptionDirectoty")
	fileName := util.GetConfigString("Nginx.subscriptionFileName")

	vmesses := make([]string, 0)

	for i := 0; i < int(vmessNum); i++ {
		cfDomaion := cdn_cf.CreateDistribution(V2rayDomainName)
		vmesses = append(vmesses, v2ray.CreateVmessFromVmess(cfDomaion, oldVmess))
	}

	subscriptionContent := v2ray.VmessToSubscription(vmesses)

	util.WriteStringToFile(dir, fileName, subscriptionContent)

}

func disableDistribution() {
	hour := util.GetConfigInt("Cloudfront.disableBeforeHours")

	distributionIDs := cdn_cf.GetDistributionsCreatedBefore(hour)

	for _, id := range distributionIDs {
		cdn_cf.DisableDistribution(id)
	}
}

func deleteDistribution() {
	hour := util.GetConfigInt("Cloudfront.deleteBeforeHours")

	distributionIDs := cdn_cf.GetDistributionsCreatedBefore(hour)

	for _, id := range distributionIDs {
		cdn_cf.DeleteDistribution(id)
	}

}
