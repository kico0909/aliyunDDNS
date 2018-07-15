package ip

import (
	"net/http"
	"io/ioutil"
	"strings"
)

func External()string{
	resp, err := http.Get("http://ip.cip.cc/")
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	content, _ := ioutil.ReadAll(resp.Body)
	return strings.Replace(string(content),"\n", "", -1)
}
