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

	for _, info := range conf.DdnsConf {
		// 获得域名解析信息
		res := domain.DomainRecordsInfo(info.Domain)

		// 查找指定域名的解析ID
		recordID := ""
		for _,v := range res.DomainRecords.Record {
			if v.RR == info.RR {
				if v.Value != IP {
					recordID = v.RecordId
				}else{
					log.Println("域名[" + info.RR + "." + info.Domain+"]指向IP未改变,无需重新解析!")
					goto STOPTHIS
				}
				break
			}
		}

		if len(recordID)>0{
			// 修改域名解析记录
			if domain.UpdateDomainRecord(info.RR, recordID, IP) {
				log.Println("[" + info.RR + "." + info.Domain+"] 更新成功!")
			}else{
				log.Println("[" + info.RR + "." + info.Domain+"] 更新失败!")
			}
		}else{

			log.Println("[" + info.RR + "." + info.Domain+"] 的解析记录未找到!")
		}
		STOPTHIS:
	}



}

type ConfigType struct {
	Appid string
	Appsecret string
	DdnsConf []configTypeChip
}
type configTypeChip struct {
	Domain string
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
