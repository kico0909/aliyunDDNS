package domain

import (
	"time"
	"log"
	"strconv"
	"net/http"
	"io/ioutil"
	"ChunkLib/codeHandler"
	"sort"
	"strings"
	"net/url"
	"encoding/json"
)

const host string = "http://alidns.aliyuncs.com";

type aliyunDns struct {
	Appid string
	Appsecrt string
}



func New(appid, appsecrt string) aliyunDns {
	return aliyunDns{appid,appsecrt}
}


// 获得指定域名的信息
type Aliyun_Domain_Info struct {
	AliDomain bool
	AvailableTtls map[string][]string
	DnsServers map[string][]string
	DomainId string
	DomainName string
	InstanceId string
	PunyCode string
	RecordLines  map[string][]string
	RequestId string
	VersionCode string
	VersionName string
}
func (_self *aliyunDns) Info(domain string) Aliyun_Domain_Info{
	var res Aliyun_Domain_Info
	params := makePublicArguments( _self.Appid )

	// 自定义参数
	params["Action"] = "DescribeDomainInfo"
	params["DomainName"] = domain

	arrKeys := sortKeys(params)

	path := makeUrl(arrKeys, params)
	sign := makeSign(path, _self.Appsecrt)
	path = path + "&Signature=" + url.QueryEscape(sign)

	url := host + "?" +path
	resq, err := http.Get(url)
	if err != nil {
		log.Print(err)
		return res
	}
	defer resq.Body.Close()

	content, _ := ioutil.ReadAll(resq.Body)



	json.Unmarshal(content, &res)

	return res
}

// 获得指定域名的解析信息
type Aliyun_Domain_Records_Info struct {
	DomainRecords struct{
		Record []aliyun_Domain_Records_Info_RecordList
	}
	PageNumber int64
	PageSize int64
}
type aliyun_Domain_Records_Info_RecordList struct {
	DomainName string
	Line string
	Locked bool
	RR string
	RecordId string
	Status string
	TTL int64
	Type string
	Value string
	Weight int64
}
func (_self *aliyunDns) DomainRecordsInfo(domain string)Aliyun_Domain_Records_Info{
	var res Aliyun_Domain_Records_Info
	params := makePublicArguments( _self.Appid )

	// 自定义参数
	params["Action"] = "DescribeDomainRecords"
	params["DomainName"] = domain
	params["PageNumber"] = "1"
	params["PageSize"] = "100"

	arrKeys := sortKeys(params)

	path := makeUrl(arrKeys, params)
	sign := makeSign(path, _self.Appsecrt)
	path = path + "&Signature=" + url.QueryEscape(sign)

	url := host + "?" +path
	resq, err := http.Get(url)
	if err != nil {
		log.Print(err)
		//return
		return res
	}
	defer resq.Body.Close()

	content, _ := ioutil.ReadAll(resq.Body)

	json.Unmarshal(content, &res)

	return res
}

// 更新指定域名的解析信息
type updateDomainRecordType struct {
	RecordId string
	RequestId string
}
func (_self *aliyunDns) UpdateDomainRecord(RR,RecordId,IP string) bool{
	var res updateDomainRecordType
	params := makePublicArguments( _self.Appid )

	// 自定义参数
	params["Action"] = "UpdateDomainRecord"
	params["RecordId"] = RecordId
	params["RR"] = RR
	params["Type"] = "A"
	params["Value"] = IP

	arrKeys := sortKeys(params)

	path := makeUrl(arrKeys, params)
	sign := makeSign(path, _self.Appsecrt)
	path = path + "&Signature=" + url.QueryEscape(sign)

	url := host + "?" +path
	resq, err := http.Get(url)
	if err != nil {
		log.Panicln(err)
	}
	defer resq.Body.Close()

	content, _ := ioutil.ReadAll(resq.Body)

	json.Unmarshal(content, &res)

	if len(res.RecordId) <1 {
		log.Println(string(content))
		return false
	}
	return true
}


// 生成公共参数
func makePublicArguments(appid string)map[string]string{

	Timestamp := time.Now().UTC().Format("2006-01-02T15:04:05Z")
	randInt := makeSignatureNonce()

	res := make(map[string]string)
	res["Format"] = "JSON"
	res["Version"] = "2015-01-09"
	res["AccessKeyId"] = appid
	res["SignatureMethod"] = "HMAC-SHA1"
	res["Timestamp"] = Timestamp
	res["SignatureVersion"] = "1.0"
	res["SignatureNonce"] = randInt

	return res
}

// 参数排序
func sortKeys(params map[string]string)[]string{
	arrKeys := make([]string, 0)
	for k,_ := range params {
		arrKeys = append(arrKeys, k)
	}
 	sort.Strings(arrKeys)

	return arrKeys
}

// 生成无序子串
func makeSignatureNonce ()string{
	return codeHandler.MD5( strconv.FormatInt(time.Now().UnixNano(), 10))
}

// 制作GETUrl
func makeUrl(keys []string, params map[string]string) string {

	for i,v := range keys {
		keys[i] += "=" + url.QueryEscape(params[v])
	}
	return strings.Join(keys, "&")

}

// 生成SIGN
func makeSign (path, secrt string)string{
	StringToSign := "GET&"+ url.QueryEscape("/")+"&" + url.QueryEscape( path)
	return codeHandler.HMAC_SHA1(StringToSign, secrt+"&")
}