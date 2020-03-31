package taoke

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/bitly/go-simplejson"

	"github.com/dshechao/go-taoke/cache"
	"github.com/nilorg/sdk/convert"
)

var (
	/**淘宝平台信息**/
	// AppKey 应用Key
	AppKeyTaobao string
	// AppSecret 秘密
	AppSecretTaobao string
	// Router 环境请求地址
	RouterTaobao string
	//API Version ...
	VersionTaobao = "2.0"
	// Session 用户登录授权成功后，TOP颁发给应用的授权信息。当此API的标签上注明：“需要授权”，则此参数必传；“不需要授权”，则此参数不需要传；“可选授权”，则此参数为可选
	Session string
	// Timeout ...
	Timeout time.Duration
	// CacheExpiration 缓存过期时间
	CacheExpiration = time.Hour
	// GetCache 获取缓存
	GetCache cache.GetCacheFunc
	// SetCache 设置缓存
	SetCache cache.SetCacheFunc

	/**京东平台信息**/
	// AppKey 应用Key
	AppKeyJingDong string
	// AppSecret 秘密
	AppSecretJingDong string
	// Router 环境请求地址
	RouterJingDong string
	//API Version ...
	VersionJingDong = "1.0"

	/**拼多多平台信息**/
	// AppKey 应用Key
	ClientIdPDD string
	// AppSecret 秘密
	ClientSecretPDD string
	// Router 环境请求地址
	RouterPDD string
	//API Version ...
	VersionPDD = "1.0"

	/**考拉海购平台信息**/
	// AppKey 应用Key
	UnionIdKL string
	// AppSecret 秘密
	SecretKL string
	// Router 环境请求地址
	RouterKL string
	//API Version ...
	VersionKL = "1.0"
)

// Parameter 参数
type Parameter map[string]interface{}

// copyParameter 复制参数
func copyParameter(srcParams Parameter) Parameter {
	newParams := make(Parameter)
	for key, value := range srcParams {
		newParams[key] = value
	}
	return newParams
}

// newCacheKey 创建缓存Key
func newCacheKey(params Parameter) string {
	cpParams := copyParameter(params)
	delete(cpParams, "session")
	delete(cpParams, "timestamp")
	delete(cpParams, "sign")

	keys := []string{}
	for k := range cpParams {
		keys = append(keys, k)
	}
	// 排序asc
	sort.Strings(keys)
	// 把所有参数名和参数值串在一起
	cacheKeyBuf := new(bytes.Buffer)
	for _, k := range keys {
		cacheKeyBuf.WriteString(k)
		cacheKeyBuf.WriteString("=")
		cacheKeyBuf.WriteString(interfaceToString(cpParams[k]))
	}
	h := md5.New()
	io.Copy(h, cacheKeyBuf)
	return hex.EncodeToString(h.Sum(nil))
}

// execute 执行API接口
func execute(param Parameter, router string) (bytes []byte, err error) {
	err = checkConfig()
	if err != nil {
		return
	}

	var req *http.Request
	req, err = http.NewRequest("POST", router, strings.NewReader(param.getRequestData()))
	if err != nil {
		return
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded;charset=utf-8")
	httpClient := &http.Client{}
	httpClient.Timeout = Timeout
	var response *http.Response
	response, err = httpClient.Do(req)
	if err != nil {
		return
	}

	if response.StatusCode != 200 {
		err = fmt.Errorf("请求错误:%d", response.StatusCode)
		return
	}
	defer response.Body.Close()
	bytes, err = ioutil.ReadAll(response.Body)
	return
}

// Execute 执行API接口
func Execute(method string, param Parameter) (res *simplejson.Json, err error) {
	p, r := setRequestData(param, method)

	var bodyBytes []byte
	bodyBytes, err = execute(p, r)
	if err != nil {
		return
	}

	return bytesToResult(bodyBytes)
}

func bytesToResult(bytes []byte) (res *simplejson.Json, err error) {
	res, err = simplejson.NewJson(bytes)
	if err != nil {
		return
	}
	if responseError, ok := res.CheckGet("error_response"); ok {
		if subMsg, subOk := responseError.CheckGet("sub_msg"); subOk {
			err = errors.New(subMsg.MustString())
		} else if zhDesc, descOk := responseError.CheckGet("zh_desc"); descOk {
			err = errors.New(zhDesc.MustString())
		} else {
			err = errors.New(responseError.Get("msg").MustString())
		}
		res = nil
	}
	return
}

// ExecuteCache 执行API接口，缓存
func ExecuteCache(method string, param Parameter) (res *simplejson.Json, err error) {
	p, r := setRequestData(param, method)

	cacheKey := newCacheKey(p)
	// 获取缓存
	if GetCache != nil {
		cacheBytes := GetCache(cacheKey)
		if len(cacheBytes) > 0 {
			res, err = simplejson.NewJson(cacheBytes)
			if err == nil && res != nil {
				return
			}
		}
	}

	var bodyBytes []byte
	bodyBytes, err = execute(param, r)
	if err != nil {
		return
	}
	res, err = bytesToResult(bodyBytes)
	if err != nil {
		return
	}
	ejsonBody, _ := res.MarshalJSON()
	// 设置缓存
	if SetCache != nil {
		go SetCache(cacheKey, ejsonBody, CacheExpiration)
	}
	return
}

// 检查配置
func checkConfig() error {
	if AppKeyTaobao == "" && AppKeyJingDong == "" && UnionIdKL == "" && ClientIdPDD == "" {
		return errors.New("至少需要设置一个平台参数")
	}
	return nil
}

//组装参数及添加公共参数
func setRequestData(p Parameter, method string) (Parameter, string) {
	platform := strings.Split(method, ".")[0]
	router := ""
	hh, _ := time.ParseDuration("8h")
	loc := time.Now().UTC().Add(hh)
	if platform == "taobao" {
		//淘宝
		p["timestamp"] = strconv.FormatInt(loc.Unix(), 10)
		p["partner_id"] = "Blant"
		p["app_key"] = AppKeyTaobao
		p["v"] = VersionTaobao
		if Session != "" {
			p["session"] = Session
		}
		p["method"] = method
		p["format"] = "json"
		p["sign_method"] = "md5"
		// 设置签名
		p["sign"] = getSign(p, AppSecretTaobao)
		router = RouterTaobao
	} else if platform == "jd" {
		//京东
		param := p
		p = Parameter{}
		p["param_json"] = param
		p["app_key"] = AppKeyJingDong
		p["v"] = VersionJingDong
		p["timestamp"] = loc.Format("2006-01-02 15:04:05")
		p["method"] = method
		p["format"] = "json"
		p["sign_method"] = "md5"
		// 设置签名
		p["sign"] = getSign(p, AppSecretJingDong)
		router = RouterJingDong
	} else if platform == "pdd" {
		//拼多多
		p["type"] = method
		p["data_type"] = "json"
		p["version"] = VersionPDD
		p["client_id"] = ClientIdPDD
		p["timestamp"] = strconv.FormatInt(loc.Unix(), 10)
		// 设置签名
		p["sign"] = getSign(p, ClientSecretPDD)
		router = RouterPDD
	} else if platform == "kaola" {
		//考拉海购
		p["method"] = method
		p["v"] = VersionKL
		p["signMethod"] = "md5"
		p["unionId"] = UnionIdKL
		p["timestamp"] = loc.Format("2006-01-02 15:04:05")
		// 设置签名
		p["sign"] = getSign(p, SecretKL)
		router = RouterKL
	}

	return p, router
}

// 获取请求数据
func (p Parameter) getRequestData() string {
	// 公共参数
	args := url.Values{}
	// 请求参数
	for key, val := range p {
		args.Set(key, interfaceToString(val))
	}
	return args.Encode()
}

// 获取签名
func getSign(params Parameter, secret string) string {
	// 获取Key
	var keys []string
	for k := range params {
		keys = append(keys, k)
	}
	// 排序asc
	sort.Strings(keys)
	// 把所有参数名和参数值串在一起
	query := bytes.NewBufferString(secret)
	for _, k := range keys {
		query.WriteString(k)
		query.WriteString(interfaceToString(params[k]))
	}
	query.WriteString(secret)
	// 使用MD5加密
	h := md5.New()
	_, _ = io.Copy(h, query)
	// 把二进制转化为大写的十六进制
	return strings.ToUpper(hex.EncodeToString(h.Sum(nil)))
}

func interfaceToString(src interface{}) string {
	if src == nil {
		panic(ErrTypeIsNil)
	}
	switch src.(type) {
	case string:
		return src.(string)
	case int, int8, int32, int64:
	case uint8, uint16, uint32, uint64:
	case float32, float64:
		return convert.ToString(src)
	}
	data, err := json.Marshal(src)
	if err != nil {
		panic(err)
	}
	return string(data)
}
