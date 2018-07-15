package main

import (
	"log"
	"ip"
	"domain"
	"config"
)

const (
	appsecrt string = "gcTAoNHXUpjjB4eKb8ZzOhjI6mdfl9"
	appid string = "LTAIP0T47hNPvDVz"
)

func main() {
	// 创建阿里云解析类
	domain := domain.New(appid, appsecrt)

	// 获得配置信息
	conf := config.Get()

	// 获得外网IP
	IP := ip.External()

	// 获得域名解析信息
	res := domain.DomainRecordsInfo("chunkding.com")

	// 查找指定域名的解析ID
	recordID := ""
	for _,v := range res.DomainRecords.Record {
		if v.RR == conf[0].RR {
			if v.Value != IP {
				recordID = v.RecordId
			}else{
				log.Println("域名指向IP未改变,无需重新解析!")
			}
			break
		}
	}

	if len(recordID)>0{
		// 修改域名解析记录
		if domain.UpdateDomainRecord(conf[0].RR, recordID, IP) {
			log.Println("dns update success!")
		}else{
			log.Println("dns update fail!")
		}
	}








	//log.Println(res)
	//log.Println(os.Args)
	//log.Println(os.Getpid())
	

}
