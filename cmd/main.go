package main

import (
	"github.com/dshechao/go-taoke"
	"log"
	"os"
	"time"
)

func init() {
	taoke.AppKey = os.Getenv("OPEN_TAOBAO_APPKEY")
	taoke.AppSecret = os.Getenv("OPEN_TAOBAO_APPSECRET")
	taoke.Router = "https://router.jd.com/api"
	taoke.V = "1.0"
	taoke.Platform = "2"
	taoke.GetCache = func(cacheKey string) []byte {
		return nil
	}
	taoke.SetCache = func(cacheKey string, value []byte, expiration time.Duration) bool {
		return true
	}
}

func main() {

	result, err := taoke.Execute("jd.union.open.goods.jingfen.query", taoke.Parameter{
		"goodsReq": taoke.Parameter{"eliteId": "1"},
	})

	if err != nil {
		log.Printf("execute error:%s\n", err)
		return
	}
	data, _ := result.MarshalJSON()
	log.Printf("result:%s\n", data)
}
