package main

import (
	"gitee.com/vblant/go-taoke"
	"log"
	"os"
	"time"
)

func init() {
	taoke.AppKey = os.Getenv("OPEN_TAOBAO_APPKEY")
	taoke.AppSecret = os.Getenv("OPEN_TAOBAO_APPSECRET")
	taoke.GetCache = func(cacheKey string) []byte {
		return nil
	}
	taoke.SetCache = func(cacheKey string, value []byte, expiration time.Duration) bool {
		return true
	}
}

func main() {

	result, err := taoke.Execute("taobao.tbk.relation.refund", taoke.Parameter{
		"search_option": map[string]interface{}{
			"page_size":   1,
			"search_type": 4, // 1-维权发起时间，2-订单结算时间（正向订单），3-维权完成时间，4-订单创建时间
			"refund_type": 1, // 1 表示2方，2表示3方
			"start_time":  "2019-07-08 00:00:00",
			"page_no":     1,
			"biz_type":    1,
		},
	})

	if err != nil {
		log.Printf("execute error:%s\n", err)
		return
	}
	data, _ := result.MarshalJSON()
	log.Printf("result:%s\n", data)
}
