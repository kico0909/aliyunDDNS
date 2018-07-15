package main

import (
	"log"
	"ip"
	"domain"
	"io/ioutil"
	"encoding/json"
)

func main() {
	// 获得配置信息
	conf := configGet()

	// 创建阿里云解析类
	domain := domain.New(conf.Appid, conf.Appsecret)

	// 获得外网IP
	IP := ip.External()

	// 获得域名解析信息
	res := domain.DomainRecordsInfo("chunkding.com")

	// 查找指定域名的解析ID
	recordID := ""
	for _,v := range res.DomainRecords.Record {
		if v.RR == conf.DdnsProfile[0].RR {
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
		if domain.UpdateDomainRecord(conf.DdnsProfile[0].RR, recordID, IP) {
			log.Println("dns update success!")
		}else{
			log.Println("dns update fail!")
		}
	}

}

type ConfigType struct {
	Appid string
	Appsecret string
	DdnsProfile []configTypeChip
}
type configTypeChip struct {
	DDNSDomain string
	RR string
}
func configGet()ConfigType{
	var res ConfigType
	fc, err := ioutil.ReadFile("./config.json")

	if err != nil {
		log.Panicln(err)
	}

	json.Unmarshal(fc, &res)

	return res

}
