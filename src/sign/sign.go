package sign

import (

	"sort"
	"log"
)


type Test struct{
	Action string
	DomainName string
	Version string
	AccessKeyId string
	SignatureVersion string
	SignatureNonce string
	Timestamp string
	PageNumber string
	PageSize string
	RRKeyWord string
	TypeKeyWord string
	ValueKeyWord string
	Format string
}

func createSign(){

}


func Make (params map[string]string ,ts, nonce, appkey string) string {
	paramsKey := make([]string,6)
	for k,_ := range params {
		paramsKey = append(paramsKey,k)
	}
	sort.Strings(paramsKey)

	StringToSign := "GET&"



	log.Println(StringToSign)
	log.Println(paramsKey)
	return ""
}





