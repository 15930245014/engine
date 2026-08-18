package dsl

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"
	"time"

	"gitlab.shudieds.com/mec/lib/utils/httpclient"
	"gitlab.shudieds.com/mec/lib/utils/jsoniter"

	"github.com/peterbourgon/diskv"
	"gitlab.shudieds.com/mec/lib/utils/convert"
	"gitlab.shudieds.com/mec/lib/utils/encrypt"
	"gitlab.shudieds.com/mec/lib/utils/maphelper"
	"gitlab.shudieds.com/zxh/engine/conf"
	"gitlab.shudieds.com/zxh/engine/types"
)

const transformBlockSize = 2

type TokenInfo struct {
	Errno   int         `json:"errno"`
	Errmsg  string      `json:"errmsg"`
	Data    Data        `json:"data"`
	TraceID string      `json:"trace_id"`
	Stack   interface{} `json:"stack"`
}

// 赛狐tocken
type SellFoxToken struct {
	Code      int         `json:"code"`
	Msg       string      `json:"msg"`
	Data      SellFoxData `json:"data"`
	RequestId string      `json:"requestId"`
}

// 赛狐tocken
type SellFoxData struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int64  `json:"expires_in"`
	GetTokenTime int64  `json:"getTokenTime"`
}

type Data struct {
	AccessToken  string `json:"accessToken"`
	ExpiresIn    int    `json:"expiresIn"`
	ExpiresOut   int    `json:"expiresOut"`
	GetTokenTime int64  `json:"getTokenTime"`
}

func (a *AST) getJJToken(expr ...types.ExprAST) (interface{}, error) {

	url, err := a.ExprASTResult(expr[0])
	if err != nil {
		return nil, err
	}

	appId, err := a.ExprASTResult(expr[1])

	if err != nil {
		return nil, err
	}

	appKey, err := a.ExprASTResult(expr[2])

	if err != nil {
		return nil, err
	}

	url, ok := convert.String(url)

	if !ok {
		return nil, errors.New("ExprASTResult")
	}

	appId, ok = convert.String(appId)

	if !ok {
		return nil, errors.New("ExprASTResult")
	}

	appKey, ok = convert.String(appKey)

	if !ok {
		return nil, errors.New("ExprASTResult")
	}

	key := encrypt.Md5sum("token")
	header := map[string]interface{}{}
	body := map[string]interface{}{}
	params := map[string]interface{}{
		"appId":  appId,
		"appKey": appKey,
	}
	d := diskv.New(diskv.Options{
		BasePath:     "data/jj",
		Transform:    blockTransform,
		CacheSizeMax: 1024 * 1024, // 1MB
	})

	token, _ := d.Read(key)

	isExpires := true

	tokenInfo := &SellFoxToken{}

	if len(token) > 0 {

		if err := json.Unmarshal(token, &tokenInfo); err != nil {
			return nil, err
		}
		if int64(tokenInfo.Data.ExpiresIn)+int64(tokenInfo.Data.ExpiresIn) > time.Now().Unix() {
			isExpires = false
		}
	}

	if isExpires {
		res, err := httpclient.DoRequest(url.(string), "GET", "", header, params, body, 10)
		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal([]byte(res), &tokenInfo); err != nil {
			return nil, err
		}

		tokenInfo.Data.GetTokenTime = time.Now().Unix()

		data, _ := json.Marshal(tokenInfo)

		d.Write(key, data)

		token = data

	}

	mp := make(map[string]interface{})

	if err := json.Unmarshal(token, &mp); err != nil {
		return nil, err
	}

	val, ok := maphelper.JDGet(mp, []string{"data", "accessToken"})
	if !ok {
		return nil, errors.New("get data err")
	}

	return val, nil

}

func (a *AST) getWDTAuth(expr ...types.ExprAST) (interface{}, error) {
	appkey, err := a.ExprASTResult(expr[0])
	if err != nil {
		return nil, err
	}

	sid, err := a.ExprASTResult(expr[1])

	if err != nil {
		return nil, err
	}

	appName, err := a.ExprASTResult(expr[2])

	if err != nil {
		return nil, err
	}

	body, err := a.ExprASTResult(expr[3])
	if err != nil {
		return nil, err
	}

	//time
	t, err := a.ExprASTResult(expr[4])
	if err != nil {
		if err != nil {
			return nil, err
		}
	}
	tStr, _ := convert.String(t)

	// body = `{
	// 	"id_list": [],
	// 	"sku_list": [],
	// 	"page_size": 100,
	// 	"page_no": 1,
	// 	"start_time": "2022-05-30 11:09:57",
	// 	"end_time": "2022-06-30 11:09:57",
	// 	"status": 1
	// }`

	m := map[string]string{
		"sid":       sid.(string),
		"appName":   appName.(string), //nJ2gK4bI3mV5bF2z
		"timestamp": tStr,
		//"body":      body.(string),
		"body": body.(string),
	}

	//获取keys
	var keys []string
	for key := range m {
		keys = append(keys, key)
	}

	//key排序
	sort.Strings(keys)

	//声明返回值
	result := make(map[string]string)
	for i := 0; i < len(keys); i++ {
		result[keys[i]] = m[keys[i]]
	}

	return encrypt.Md5sum(linkParams(result, appkey.(string))), nil
}

func (a *AST) getJSTsign(expr ...types.ExprAST) (interface{}, error) {

	accessToken, err := a.ExprASTResult(expr[0])
	if err != nil {
		return nil, err
	}

	version, err := a.ExprASTResult(expr[1])
	if err != nil {
		return nil, err
	}

	appKey, err := a.ExprASTResult(expr[2])
	if err != nil {
		return nil, err
	}

	biz, err := a.ExprASTResult(expr[3])
	if err != nil {
		return nil, err
	}

	appSecret, err := a.ExprASTResult(expr[4])
	if err != nil {
		return nil, err
	}

	//m := make(map[string]string)
	m := map[string]string{
		"charset":      "utf-8",
		"access_token": accessToken.(string),
		"timestamp":    strconv.FormatInt(time.Now().Unix(), 10),
		"version":      version.(string),
		"app_key":      appKey.(string),
		"biz":          biz.(string),
	}

	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	resultStr := ""
	for _, k := range keys {
		resultStr = fmt.Sprintf("%s%s%s", resultStr, k, m[k])
	}

	resultStr = appSecret.(string) + resultStr
	hash := md5.Sum([]byte(resultStr))
	// 将哈希值转换为十六进制字符串
	hexHash := hex.EncodeToString(hash[:])

	return hexHash, nil
}

func (a *AST) getJDYSessionId(expr ...types.ExprAST) (interface{}, error) {

	header := map[string]interface{}{
		"Content-Type": "application/json",
	}
	body := map[string]interface{}{
		"acctid":   "20161107155638",
		"username": "admin10",
		"password": "xwf@123456",
		"lcid":     2052,
	}
	params := map[string]interface{}{}

	url := "https://worthfind.ik3cloud.com/k3cloud/Kingdee.BOS.WebApi.ServicesStub.AuthService.ValidateUser.common.kdsvc"
	res, err := httpclient.DoRequest(url, "POST", "", header, params, body, 10)
	if err != nil {
		return nil, err
	}

	//tokenInfo := JDYGenerated{}
	response := make(map[string]interface{})

	if err := json.Unmarshal([]byte(res), &response); err != nil {
		return nil, err
	}

	val, ok := maphelper.JDGet(response, []string{"KDSVCSessionId"})
	if !ok {
		return nil, errors.New("get data err")
	}

	return val, nil
}

// 数据存储到文件的转换
func blockTransform(s string) []string {
	var (
		sliceSize = len(s) / transformBlockSize
		pathSlice = make([]string, sliceSize)
	)
	for i := 0; i < sliceSize; i++ {
		from, to := i*transformBlockSize, (i*transformBlockSize)+transformBlockSize
		pathSlice[i] = s[from:to]
	}
	return pathSlice
}

// 旺店通auth 转换
func linkParams(mapData map[string]string, secret string) string {
	var sb strings.Builder
	sb.WriteString(secret)

	for k, v := range mapData {
		if k == "sign" {
			continue
		}
		sb.WriteString(k)
		sb.WriteString(v)
	}

	sb.WriteString(secret)
	return sb.String()
}

// 获取Eenecan 签名
func (a *AST) getEccangSign(expr ...types.ExprAST) (interface{}, error) {

	signStr := "app_key=%s&biz_content=%s&charset=UTF-8&interface_method=%s&nonce_str=%s&service_id=%s&sign_type=MD5&timestamp=%d&version=%s%s"
	appKey, err := a.ExprASTResult(expr[0])
	if err != nil {
		return nil, err
	}

	bizContent, err := a.ExprASTResult(expr[1])
	if err != nil {
		return nil, err
	}

	params, ok := convert.ToMap(bizContent)
	if !ok {
		return nil, errors.New("bizContent is 不合法")
	}

	bt, _ := jsoniter.MarshalV2(params)

	bizContent = convert.BytesToString(bt)

	interfaceMethod, err := a.ExprASTResult(expr[2])
	if err != nil {
		return nil, err
	}

	nonceStr, err := a.ExprASTResult(expr[3])
	if err != nil {
		return nil, err
	}

	serviceId, err := a.ExprASTResult(expr[4])
	if err != nil {
		return nil, err
	}

	timestamp, err := a.ExprASTResult(expr[5])
	if err != nil {
		return nil, err
	}
	timestamp, _ = convert.Int64(timestamp)
	version, err := a.ExprASTResult(expr[6])
	if err != nil {
		return nil, err
	}

	appSecret, err := a.ExprASTResult(expr[7])
	if err != nil {
		return nil, err
	}

	signStr = fmt.Sprintf(signStr, appKey, bizContent, interfaceMethod, nonceStr, serviceId, timestamp, version, appSecret)

	fmt.Printf("getEccangSign is signStr :%v\n", signStr)
	return encrypt.Md5sum(signStr), nil

	// header := map[string]interface{}{
	// 	"Content-Type": "application/json",
	// }
	// body := map[string]interface{}{
	// 	"app_key":          appKey,
	// 	"service_id":       serviceId,
	// 	"nonce_str":        nonceStr,
	// 	"interface_method": interfaceMethod,
	// 	"biz_content":      bizContent,
	// 	"sign":             encrypt.Md5sum(signStr),
	// 	"sign_type":        "MD5",
	// 	"charset":          "UTF-8",
	// 	"timestamp":        timestamp,
	// 	"version":          version,
	// }

	// ret, err := httpclient.DoRequest("http://openapi-web.eccang.com/openApi/api/unity", "POST", "application/json", header, params, body, 10)
	// //ret, err := httpclient.DoRequest("http://testnew-openapi-web.eccang.com", "POST", "application/json", header, params, body, 10)
	// return ret, err
}

// 获取赛狐token

func (a *AST) getSellfoxToken(expr ...types.ExprAST) (interface{}, error) {
	clientId, err := a.ExprASTResult(expr[0])
	if err != nil {
		return nil, err
	}
	clientSecret, err := a.ExprASTResult(expr[1])
	if err != nil {
		return nil, err
	}

	header := map[string]interface{}{
		"Content-Type": "application/json",
	}
	body := map[string]interface{}{}
	params := map[string]interface{}{
		"client_id":     clientId,
		"client_secret": clientSecret,
		"grant_type":    "client_credentials",
	}

	//ret, err := httpclient.DoRequest("https://openapi.sellfox.com/api/oauth/v2/token.json", "GET", "application/json", header, params, body, 10)
	//return ret, err

	key := encrypt.Md5sum("sellfoxtoken")
	d := diskv.New(diskv.Options{
		BasePath: "data/sellfox",
		//Transform:    blockTransform,
		CacheSizeMax: 1024 * 1024, // 1MB
	})

	//缓存读到 byte
	token, _ := d.Read(key)

	//是否无效
	isNotExpires := true

	//解析后
	tokenInfo := &SellFoxToken{}

	if len(token) > 0 {
		if err := json.Unmarshal(token, &tokenInfo); err != nil {
			return nil, err
		}
		//获取token时间 + 过期时间  12：00：02   12：00：01
		if tokenInfo.Data.GetTokenTime+tokenInfo.Data.ExpiresIn > time.Now().Unix() {
			isNotExpires = false
		}
	}

	if isNotExpires == true {
		res, err := httpclient.DoRequest("https://openapi.sellfox.com/api/oauth/v2/token.json", "GET", "application/json", header, params, body, 10)
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal([]byte(res), &tokenInfo); err != nil {
			return nil, err
		}

		tokenInfo.Data.GetTokenTime = time.Now().Unix()

		data, _ := json.Marshal(tokenInfo)
		_ = d.Write(key, data)
	}
	relToken := tokenInfo.Data.AccessToken
	return relToken, nil
}

// 获取赛狐sign
// token, clientId, method,url, secret interface{}, timestamp string
func (a *AST) getSellfoxSign(expr ...types.ExprAST) (interface{}, error) {
	tokenVal, err := a.ExprASTResult(expr[0])
	if err != nil {
		return nil, errors.New("getSellfoxSign 获取token错误：" + err.Error())
	}
	token, ok := convert.String(tokenVal)
	if !ok {
		return nil, errors.New("getSellfoxSign 获取token错误：类型应为str")
	}

	clientIdVal, err := a.ExprASTResult(expr[1])
	if err != nil {
		return nil, errors.New("getSellfoxSign 获取clientId错误：" + err.Error())
	}
	clientId, ok := convert.String(clientIdVal)
	if !ok {
		return nil, errors.New("getSellfoxSign 获取clientId错误：类型应为str")
	}

	methodVal, err := a.ExprASTResult(expr[2])
	if err != nil {
		return nil, errors.New("getSellfoxSign 获取method错误：" + err.Error())
	}
	method, ok := convert.String(methodVal)
	if !ok {
		return nil, errors.New("getSellfoxSign 获取method错误：类型应为str")
	}

	urlVal, err := a.ExprASTResult(expr[3])
	if err != nil {
		return nil, errors.New("getSellfoxSign 获取url错误：" + err.Error())
	}
	url1, ok := convert.String(urlVal)
	if !ok {
		return nil, errors.New("getSellfoxSign 获取url错误：类型应为str")
	}

	secretVal, err := a.ExprASTResult(expr[4])
	if err != nil {
		return nil, errors.New("getSellfoxSign 获取secrets错误：" + err.Error())
	}
	secret, ok := convert.String(secretVal)
	if !ok {
		return nil, errors.New("getSellfoxSign 获取secret错误：类型应为str")
	}

	timestampVal, err := a.ExprASTResult(expr[5])
	if err != nil {
		return nil, errors.New("getSellfoxSign 获取timestamp错误：" + err.Error())
	}

	timestamp, ok := convert.String(timestampVal)
	if !ok {
		return nil, errors.New("getSellfoxSign 获取timestamp错误：类型不能转换成str")
	}

	nonceVal, err := a.ExprASTResult(expr[6])
	if err != nil {
		return nil, errors.New("getSellfoxSign 获取随机数错误：" + err.Error())
	}

	nonce, ok := convert.Int(nonceVal)
	if !ok {
		return nil, errors.New("getSellfoxSign 获取timestamp错误：类型不能转换成int")
	}

	params := map[string]string{
		"access_token": token,
		"client_id":    clientId,
		"method":       method,
		"nonce":        strconv.Itoa(nonce),
		"timestamp":    timestamp,
		"url":          url1,
	}

	// 参数排序
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var data strings.Builder
	for i, k := range keys {
		if i > 0 {
			data.WriteString("&")
		}
		data.WriteString(k)
		data.WriteString("=")
		data.WriteString(params[k])
	}

	sign, err := hmacsha256(secret, data.String())
	if err != nil {
		return "", errors.New("getSellfoxSign 获取签名错误:" + err.Error())
	}
	return sign, nil

}

// TODO
// 测试获取赛狐数据| 删除
// func getData(token, client_id, client_secret, timestamp, nonce, sign, url string) {
func getData(sign, url string) {

	url = "https://openapi.sellfox.com/api/order/pageList.json?" + url + "&sign=" + sign

	url = strings.Replace(url, "&url=/api/order/pageList.json", "", 1)
	url = strings.Replace(url, "&method=post", "", 1)

	header := map[string]interface{}{
		//"Content-Type": "application/json",
	}
	body := map[string]interface{}{}
	//intNonce, _ := strconv.ParseInt(nonce, 10, 64)

	params := map[string]interface{}{}

	fmt.Println(params)

	ret, err := httpclient.DoRequest(url, "POST", "application/json", header, params, body, 10)
	fmt.Println(ret, err)
}

func hmacsha256(key, data string) (string, error) {
	h := hmac.New(sha256.New, []byte(key))
	_, err := h.Write([]byte(data))
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

// TODO
//测试
// getWDTSign 获取旺店通

func (a *AST) getWDTQMSign(expr ...types.ExprAST) (interface{}, error) {
	inputParams, err := a.ExprASTResult(expr[0])
	if err != nil {
		return nil, err
	}

	inputParamsMap, ok := convert.ToMap(inputParams)
	if !ok {
		return nil, errors.New(" getWDTSign  的参数 不是一个 map")
	}

	//url := "http://hu3cgwt0tc.api.taobao.com/router/qm"
	appsecret, _ := convert.String(inputParamsMap["appsecret"])
	if len(appsecret) <= 0 {
		return nil, errors.New("getWDTSign appsecret is null")
	}

	delete(inputParamsMap, "appsecret")
	//delete(inputParamsMap, "pager")
	//delete(inputParamsMap, "params")

	inputParamsMap["sign_method"] = "md5"
	inputParamsMap["format"] = "json"

	keys := make([]string, 0, len(inputParamsMap))
	for k := range inputParamsMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// 奇门签名
	sign, err := signTopRequest(inputParamsMap, appsecret, "md5")
	if err != nil {
		return nil, errors.New("getWDTSign 请求错误：" + sign)
	}
	return sign, nil

}

// pdd
func (a *AST) getPDDSign(expr ...types.ExprAST) (interface{}, error) {
	v1, err := a.ExprASTResult(expr[0])
	if err != nil {
		return nil, err
	}
	params, ok := convert.ToMap(v1)
	if !ok {
		return nil, errors.New(" getPDDSign  的参数 不是一个 map")
	}
	v2, err := a.ExprASTResult(expr[1])
	if err != nil {
		return nil, err
	}
	secret, _ := convert.String(v2)
	if len(secret) <= 0 {
		return nil, errors.New("getPDDSign  secret is null")
	}

	// 第一步：检查参数是否已经排序
	keys := make([]string, 0, len(params))
	for key := range params {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	// 第二步：把所有参数名和参数值串在一起
	var query strings.Builder
	query.WriteString(secret)
	for _, key := range keys {
		value := params[key]
		if key != "" && value != "" {
			query.WriteString(key)
			strVal, ok := convert.String(value)
			if !ok {
				return "", fmt.Errorf("%s=>%s转换错误", key, strVal)
			}
			query.WriteString(strVal)
		}
	}

	// 第三步：使用MD5/HMAC加密
	var bytes []byte
	query.WriteString(secret)
	bytes, err = encryptMD5(query.String())
	if err != nil {
		return "", err
	}

	// 第四步：把二进制转化为大写的十六进制字符串
	//return byte2hex(bytes), nil
	return strings.ToUpper(hex.EncodeToString(bytes)), nil
}
func GetPDDSign(params map[string]interface{}, secret string) (string, error) {
	var err error
	// 第一步：检查参数是否已经排序
	keys := make([]string, 0, len(params))
	for key := range params {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	// 第二步：把所有参数名和参数值串在一起
	var query strings.Builder
	query.WriteString(secret)
	for _, key := range keys {
		value := params[key]
		if key != "" && value != "" {
			query.WriteString(key)
			strVal, ok := convert.String(value)
			if !ok {
				return "", fmt.Errorf("%s=>%s转换错误", key, strVal)
			}
			query.WriteString(strVal)
		}
	}

	// 第三步：使用MD5/HMAC加密
	var bytes []byte
	query.WriteString(secret)
	bytes, err = encryptMD5(query.String())
	if err != nil {
		return "", err
	}

	// 第四步：把二进制转化为大写的十六进制字符串
	//return byte2hex(bytes), nil
	return strings.ToUpper(hex.EncodeToString(bytes)), nil
}

// 奇门
func signTopRequest(params map[string]interface{}, secret string, signMethod string) (string, error) {
	// 第一步：检查参数是否已经排序
	keys := make([]string, 0, len(params))
	for key := range params {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	// 第二步：把所有参数名和参数值串在一起
	var query strings.Builder
	if signMethod == conf.SignMethodMD5 {
		query.WriteString(secret)
	}
	for _, key := range keys {
		value := params[key]
		if key != "" && value != "" {
			query.WriteString(key)
			strVal, ok := convert.String(value)
			if !ok {
				return "", fmt.Errorf("%s=>%s转换错误", key, strVal)
			}
			query.WriteString(strVal)
		}
	}

	// 移除换行符（如果存在的话）

	// 第三步：使用MD5/HMAC加密
	var bytes []byte
	var err error
	if signMethod == conf.SignMethodHMAC {
		bytes, err = encryptHMAC(query.String(), secret)
	} else if signMethod == conf.SignMethodHMACSHA256 {
		bytes, err = encryptHMACSHA256(query.String(), secret)
	} else {
		query.WriteString(secret)
		bytes, err = encryptMD5(query.String())
	}
	if err != nil {
		return "", err
	}

	// 第四步：把二进制转化为大写的十六进制字符串
	//return byte2hex(bytes), nil
	return strings.ToUpper(hex.EncodeToString(bytes)), nil
}

// 获取旺店通 sign
func (a *AST) getWDTSign(expr ...types.ExprAST) (interface{}, error) {

	//secret, salt string, params map[string]interface{}
	inputSecret, err := a.ExprASTResult(expr[0])
	if err != nil {
		return nil, err
	}
	secret, _ := convert.String(inputSecret)

	inputSalt, err := a.ExprASTResult(expr[1])
	if err != nil {
		return nil, err
	}
	salt, _ := convert.String(inputSalt)

	inputParams, err := a.ExprASTResult(expr[2])
	if err != nil {
		return nil, err
	}

	params, ok := convert.ToMap(inputParams)
	if !ok {
		return nil, errors.New(" getWDTSign  的参数 不是一个 map")
	}

	// Step 1: 将所有参数放入一个新的 map 中
	signParams := make(map[string]interface{})
	for k, v := range params {
		if k != "wdt_sign" { // 排除 wdt_sign
			signParams[k] = v
		}
	}
	signParams["wdt_salt"] = salt

	// Step 2: 对参数进行排序并拼接
	keys := make([]string, 0, len(signParams))
	for k := range signParams {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Step 3: 拼接参数值
	var builder strings.Builder
	builder.WriteString(secret)
	for _, k := range keys {
		// 特别处理 JSON 字符串类型
		value := formatValue(signParams[k])
		builder.WriteString(k)
		builder.WriteString(value)
	}

	builder.WriteString(secret)

	// Step 4: 计算 MD5 值
	signature := fmt.Sprintf("%x", md5.Sum([]byte(builder.String())))
	return strings.ToUpper(signature), nil
}

// formatValue 根据类型格式化参数值
func formatValue(value interface{}) string {
	switch v := value.(type) {
	case string:
		return v
	case map[string]interface{}: // 如果是 map，需要递归排序
		keys := make([]string, 0, len(v))
		for k := range v {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		var builder strings.Builder
		for _, k := range keys {
			builder.WriteString(k)
			builder.WriteString(formatValue(v[k]))
		}
		return builder.String()
	case []interface{}: // 如果是列表，按顺序拼接
		var builder strings.Builder
		for _, item := range v {
			builder.WriteString(formatValue(item))
		}
		return builder.String()
	case bool: // 布尔值转为字符串
		if v {
			return "true"
		}
		return "false"
	case nil: // null 不参与签名
		return ""
	default: // 其他类型转为字符串
		return fmt.Sprintf("%v", v)
	}
}

func encryptHMAC(data string, secret string) ([]byte, error) {
	key := []byte(secret)
	h := hmac.New(md5.New, key)
	_, err := io.WriteString(h, data)
	if err != nil {
		return nil, err
	}
	return h.Sum(nil), nil
}

func encryptHMACSHA256(data string, secret string) ([]byte, error) {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	_, err := io.WriteString(h, data)
	if err != nil {
		return nil, err
	}
	return h.Sum(nil), nil
}

func encryptMD5(data string) ([]byte, error) {
	h := md5.New()
	_, err := io.WriteString(h, data)
	if err != nil {
		return nil, err
	}
	return h.Sum(nil), nil
}
