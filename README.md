#淘宝客 淘宝联盟，京东联盟基础api请求包
淘宝客,京东Api、淘宝开放平台Api,京东联盟请求基础SDK

# 淘宝API

[sign算法](http://open.taobao.com/doc.htm?docId=101617&docType=1)

[淘宝Session](https://oauth.taobao.com/authorize?response_type=token&client_id=24840730)

# Example-Taobao 
```go
package main

import (
	"fmt"
	"github.com/dshechao/go-taoke"
)

func init() {
	taoke.AppKey = ""
	taoke.AppSecret = ""
    taoke.Platform = "1"
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

# Example - JingDong
```go
package main

import (
	"github.com/dshechao/go-taoke"
	"log"
)

func init() {
	taoke.AppKey = ""
	taoke.AppSecret = ""
    taoke.Router ="https://router.jd.com/api"
	taoke.V = "1.0"
	taoke.Platform = "2"
}

func main() {

	result, err := taoke.Execute("jd.union.open.goods.jingfen.query1", taoke.Parameter{
		"goodsReq":taoke.Parameter{"eliteId":    "1"},
	})

	if err != nil {
		log.Printf("execute error:%s\n", err)
		return
	}
	data, _ := result.MarshalJSON()
	log.Printf("result:%s\n", data)
}


```
# Example - 多多进宝
```go
package main

import (
	"github.com/dshechao/go-taoke"
	"log"
)

func init() {
	taoke.AppKey = "client_id"
	taoke.AppSecret = "client_secret"
	taoke.Router = "https://gw-api.pinduoduo.com/api/router"
	taoke.Platform = "3"
    taoke.V = "1.0"
}

func main() {

	result, err := taoke.Execute("pdd.ddk.goods.search", taoke.Parameter{
    		"keyword":"华为手机",
    	})

	if err != nil {
		log.Printf("execute error:%s\n", err)
		return
	}
	data, _ := result.MarshalJSON()
	log.Printf("result:%s\n", data)
}


```
# Example - 考拉赚客
```go
package main

import (
	"github.com/dshechao/go-taoke"
	"log"
)

func init() {
    taoke.AppKey = "赚客ID"
    taoke.AppSecret = "AppSecret"
    taoke.Router = "https://cps.kaola.com/zhuanke/api"
    taoke.Platform = "4"
    taoke.V = "1.0"
}

func main() {

	result, err := taoke.Execute("kaola.zhuanke.api.queryRecommendGoodsList", taoke.Parameter{
    		"sortType":1,
    		"pageIndex":1,
    	})

	if err != nil {
		log.Printf("execute error:%s\n", err)
		return
	}
	data, _ := result.MarshalJSON()
	log.Printf("result:%s\n", data)
}


```