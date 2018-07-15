package config

import (
	"log"
	"io/ioutil"
	"encoding/json"
)

type ConfigType []configTypeChip
type configTypeChip struct {
	DDNSDomain string
	RR string
}
func Get()ConfigType{
	var res ConfigType
	fc, err := ioutil.ReadFile("config.json")

	if err != nil {
		log.Panicln(err)
	}

	json.Unmarshal(fc, &res)

	return res

}



