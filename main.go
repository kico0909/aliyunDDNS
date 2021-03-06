package main

import (
	"domain"
	"encoding/json"
	"io/ioutil"
	"ip"
	"log"
	"mylog"
	"strconv"
	"time"
)

func main() {
	// 获得配置信息
	conf := configGet()

	// 创建阿里云解析类
	domain := domain.New(conf.Appid, conf.Appsecret)

	// 设置日志的单文件尺寸
	mylog.SetLogMaxSize(conf.Log.FileMaxSize)

	// 设置日志位置
	mylog.SetLogPath(conf.Log.Path, conf.Log.FileName)

	// 启动日志
	mylog.LogStart()

	nowDate := time.Now().Format("\r\n 2006-01-02 15:04:05")

	mylog.Record(nowDate + " 触发更新!")

	// 获得外网IP
	IP := ip.External()

	mylog.Record("\t当前解析IP[" + IP + "]\r\n\t配置文件长度 [ " + strconv.FormatInt(int64(len(conf.DdnsConf)), 10) + " ]")

	for _, info := range conf.DdnsConf {

		// 获得域名解析信息
		res := domain.DomainRecordsInfo(info.Domain)

		// 查找指定域名的解析ID
		recordID := ""
		for _, v := range res.DomainRecords.Record {
			if v.RR == info.RR {
				if v.Value != IP {
					recordID = v.RecordId
				} else {
					mylog.Record("\t域名[" + info.RR + "." + info.Domain + "]指向IP未改变,无需重新解析!")
					goto STOPTHIS
				}
				break
			}
		}

		if len(recordID) > 0 {
			// 修改域名解析记录
			success := domain.UpdateDomainRecord(info.RR, recordID, IP)

			if success {
				mylog.Record("\t[" + info.RR + "." + info.Domain + "] 更新成功!")
			} else {
				mylog.Record("\t[" + info.RR + "." + info.Domain + "] 更新失败!")
			}

		} else {
			mylog.Record("\t[" + info.RR + "." + info.Domain + "] 的解析记录未找到!")
		}
	STOPTHIS:
		mylog.LogStop()
	}

}

type ConfigType struct {
	Appid     string	`json:"appid"`
	Appsecret string	`json:"appsecret"`
	Log struct{
		Path	string	`json:"path"`
		FileName string	`json:"fileName"`
		FileMaxSize int64	`json:"fileMaxSize"`
	}	`json:"log"`
	DdnsConf  []struct {
		Domain string	`json:"domain"`
		RR     string	`json:"rr"`
	}	`json:"ddnsConf"`
}

func configGet() ConfigType {

	var res ConfigType

	fc, err := ioutil.ReadFile("./config.json")

	if err != nil {
		log.Panicln(err)
	}

	json.Unmarshal(fc, &res)

	return res

}
