# taobaogo
淘宝客,京东Api、淘宝开放平台Api,京东联盟请求基础SDK

# 淘宝API

[sign算法](http://open.taobao.com/doc.htm?docId=101617&docType=1)

[淘宝Session](https://oauth.taobao.com/authorize?response_type=token&client_id=24840730)

# Example 
```go
package main

import (
	"fmt"
	 "gitee.com/vblant/go-taoke"
)

func init() {
	taoke.AppKey = ""
	taoke.AppSecret = ""
	taoke.Router = "http://gw.api.taobao.com/router/rest"
}

func main() {
	res, err := taoke.Execute("taobao.tbk.dg.material.optional", taoke.Parameter{
		"adzone_id":"",
		"q":      "华为",
		"cat":    "16,18",
	})

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("商品数量:", res.Get("tbk_item_get_response").Get("total_results").MustInt())
	var imtes []interface{}
	imtes, _ = res.Get("tbk_item_get_response").Get("results").Get("n_tbk_item").Array()
	for _, v := range imtes {
		fmt.Println("======")
		item := v.(map[string]interface{})
		fmt.Println("商品名称:", item["title"])
		fmt.Println("商品价格:", item["reserve_price"])
		fmt.Println("商品链接:", item["item_url"])
	}
}

```