package main

import (
	"github.com/dshechao/go-taoke"
	"log"
)

func init() {
	taoke.AppKeyJingDong = ""
	taoke.AppSecretJingDong = ""
	taoke.VersionJingDong = "1.0"
	taoke.RouterJingDong = "https://router.jd.com/api"

}

func main() {

	method := "jd.union.open.goods.jingfen.query"
	param := taoke.Parameter{}
	param["goodsReq"] = taoke.Parameter{
		"eliteId": 1,
	}
	result, err := taoke.Execute(method, param)

	if err != nil {
		log.Println(err)
		return
	}
	data, _ := result.MarshalJSON()
	log.Printf("result:%s\n", data)
}
