package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	define "github.com/TencentBlueKing/bk-plugin-client-go/core/define"
	"github.com/TencentBlueKing/bk-plugin-client-go/core/utils"
	"io"
	"net/http"
)

type Client struct {
	UsePluginService            bool
	PluginServiceApiGwAppCode   string
	PluginServiceApiGwAppSecret string
	ApiGwEnvironment            string
	ApiGwNetworkProtocal        string
	ApiGwUrlSuffix              string
	BkAppInvokePaasRetryNum     int
	ApiGwUserAuthKeyName        string
	PluginCode                  string
	PluginHost                  string
	PluginApiGwName             string
}

func GetClient(pluginCode string) (Client, error) {
	client := Client{
		UsePluginService:            define.UsePluginService,
		PluginServiceApiGwAppCode:   define.PluginServiceApiGwAppCode,
		PluginServiceApiGwAppSecret: define.PluginServiceApiGwAppSecret,
		ApiGwEnvironment:            define.ApiGwEnvironment,
		ApiGwNetworkProtocal:        define.ApiGwNetworkProtocal,
		ApiGwUrlSuffix:              define.ApiGwUrlSuffix,
		BkAppInvokePaasRetryNum:     utils.CovertStrInt(define.BkAppInvokePaasRetryNum),
		ApiGwUserAuthKeyName:        define.ApiGwUserAuthKeyName,
		PluginCode:                  pluginCode,
	}
	if pluginCode != "" {
		// 如果用户不传 插件 code 则不进行详情获取
		detail, _ := client.GetPluginAppDetail()
		client.PluginHost = detail.Url
		client.PluginApiGwName = detail.ApiGwName
	}
	return client, nil
}

func (client *Client) preparePassApiRequest(pathParams string) string {
	// 组装Paas的url
	url := fmt.Sprintf("%s://paasv3.%s/%s/%s", client.ApiGwNetworkProtocal, client.ApiGwUrlSuffix, client.ApiGwEnvironment, pathParams)
	return url

}

func (client *Client) prepareApigwApiRequest(pathParams string) string {
	// 组装插件网关的url
	url := fmt.Sprintf("%s://%s.%s/%s/%s", client.ApiGwNetworkProtocal, client.PluginApiGwName, client.ApiGwUrlSuffix, client.ApiGwEnvironment, pathParams)
	return url
}

func (client *Client) requestApiAndErrorRetry(url string, method string, body []byte) ([]byte, error) {
	// 如果失败则进行重试，默认重试次数为3次
	headers := map[string]interface{}{
		"bk_app_code": client.PluginServiceApiGwAppCode, "bk_app_secret": client.PluginServiceApiGwAppSecret,
	}

	authHeaders, _ := json.Marshal(headers)
	httpClient := &http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return []byte{}, fmt.Errorf("url请求失败: %s", err)
	}

	for i := 0; i < client.BkAppInvokePaasRetryNum; i++ {
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Bkapi-Authorization", string(authHeaders))
		resp, err := httpClient.Do(req)
		if err != nil {
			return []byte{}, fmt.Errorf("url请求失败: %s", err)
		}
		if resp.StatusCode == 200 {
			response, _ := io.ReadAll(resp.Body)
			err = resp.Body.Close()
			if err != nil {
				return []byte{}, err
			}
			return response, nil
		}
	}
	return []byte{}, nil
}

func (client *Client) requestApi(url string, method string, body []byte) ([]byte, error) {
	// 请求接口
	headers := map[string]interface{}{
		"bk_app_code": client.PluginServiceApiGwAppCode, "bk_app_secret": client.PluginServiceApiGwAppSecret,
	}

	authHeaders, _ := json.Marshal(headers)
	httpClient := &http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return []byte{}, fmt.Errorf("url请求失败: %s", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Bkapi-Authorization", string(authHeaders))
	resp, err := httpClient.Do(req)
	if err != nil {
		return []byte{}, fmt.Errorf("url请求失败: %s", err)
	}
	if resp.StatusCode != 200 {
		message, err := io.ReadAll(resp.Body)
		if err != nil {
			message = []byte{}
		}
		return []byte{}, fmt.Errorf("请求失败，返回值非200, StatusCode: %d, err: %s", resp.StatusCode, string(message))
	}
	response, err := io.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}
	err = resp.Body.Close()
	if err != nil {
		return []byte{}, err
	}
	return response, nil
}

func (client *Client) GetPluginAppDetail() (*define.PluginAppDetail, error) {
	// 获取某个插件的详情信息, 返回插件的host. code 以及网关地址等信息
	pathParams := fmt.Sprintf("system/bk_plugins/%s/", client.PluginCode)
	url := client.preparePassApiRequest(pathParams)
	response, err := client.requestApiAndErrorRetry(url, "GET", nil)
	if err != nil {
		return nil, err
	}
	var result interface{}
	err = json.Unmarshal(response, &result)

	if err != nil {
		return nil, err
	}
	resultMap := result.(map[string]interface{})
	pluginDetail := resultMap["plugin"].(map[string]interface{})
	pluginProfile := resultMap["profile"].(map[string]interface{})
	deployedStatuses := resultMap["deployed_statuses"].(map[string]interface{})

	envDetail := deployedStatuses["stag"].(map[string]interface{})

	if !envDetail["deployed"].(bool) {
		return nil, fmt.Errorf("模板网关没有注册")
	}

	addresses := envDetail["addresses"].([]interface{})

	var hosts []string
	var defaultHost string
	for i := 0; i < len(addresses); i++ {
		address := addresses[i].(map[string]interface{})
		addressUrl := address["address"].(string)
		addressType := int(address["type"].(float64))
		if addressType == define.DefaultHostType {
			defaultHost = fmt.Sprintf("%s/bk_plugin", addressUrl)
			hosts = append(hosts, addressUrl)
		}
	}

	pluginCode := pluginDetail["code"].(string)
	name := pluginDetail["name"].(string)
	apiGwName := pluginProfile["api_gw_name"].(string)

	return &define.PluginAppDetail{
		Url:       defaultHost,
		Urls:      hosts,
		Name:      name,
		ApiGwName: apiGwName,
		Code:      pluginCode,
	}, nil
}

func (client *Client) GetPluginDetail(version string) (*define.PluginDetail, error) {
	// 返回插件的detail 信息, 输入输出，上下文等信息
	url := fmt.Sprintf("%s/detail/%s", client.PluginHost, version)
	response, err := client.requestApiAndErrorRetry(url, "GET", nil)
	if err != nil {
		return nil, fmt.Errorf("请求失败,%s", err)
	}
	var detailData define.PluginDetailReturn
	err = json.Unmarshal(response, &detailData)
	if err != nil {
		return nil, fmt.Errorf("json反序列化失败,%s", err)
	}
	if !detailData.Result {
		return nil, fmt.Errorf("请求失败, %s", detailData.Message)
	}
	return &detailData.Data, nil
}

func (client *Client) GetPluginMeta() (*define.PluginMeta, error) {
	// 返回某个插件的元信息，版本什么的
	url := fmt.Sprintf("%s/meta/", client.PluginHost)
	response, err := client.requestApiAndErrorRetry(url, "GET", nil)
	if err != nil {
		return nil, fmt.Errorf("请求异常%s", err)
	}
	var metaData define.PluginMetaReturn
	err = json.Unmarshal(response, &metaData)
	if err != nil {
		return nil, fmt.Errorf("请求异常%s", err)
	}
	if !metaData.Result {
		return nil, fmt.Errorf("请求失败%s", metaData.Message)
	}

	return &metaData.Data, nil
}

func (client *Client) Invoke(version string, data map[string]interface{}) (*define.PluginInvokeResult, error) {
	// 执行某个插件
	pathParam := fmt.Sprintf("invoke/%s/", version)
	url := client.prepareApigwApiRequest(pathParam)
	reqBody, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("请求json序列化失败, err: %s", err)
	}
	response, err := client.requestApi(url, "POST", reqBody)
	if err != nil {
		return nil, fmt.Errorf("接口请求失败, err: %s", err)
	}
	var result define.PluginInvokeResult
	err = json.Unmarshal(response, &result)
	if err != nil {
		return nil, fmt.Errorf("json 反序列化失败, err: %s", err)
	}
	if !result.Result {
		return nil, fmt.Errorf("插件请求失败, message: %s, trace_id: %s", result.Message, result.TraceId)
	}

	return &result, nil
}

func (client *Client) GetSchedule(traceId string) (*define.PluginScheduleResult, error) {
	// 查询插件的状态
	url := fmt.Sprintf("%s/schedule/%s", client.PluginHost, traceId)
	response, err := client.requestApiAndErrorRetry(url, "GET", nil)
	if err != nil {
		return nil, fmt.Errorf("接口请求失败, err: %s", err)
	}
	var result define.PluginScheduleResult
	err = json.Unmarshal(response, &result)
	if err != nil {
		return nil, fmt.Errorf("json 反序列化失败, err: %s", err)
	}

	if !result.Result {
		return nil, fmt.Errorf("插件请求失败, message: %s, trace_id: %s", result.Message, result.TraceId)
	}

	return &result, nil

}

func (client *Client) GetPluginLogs(traceId string, scrollId interface{}) (*define.PluginLogResult, error) {
	// 查询插件的日志
	pathParam := fmt.Sprintf("system/bk_plugins/%s/logs?trace_id=%s", client.PluginCode, traceId)
	url := client.preparePassApiRequest(pathParam)
	if scrollId != nil {
		url = fmt.Sprintf("%s&scroll_id=%s", url, scrollId.(string))
	}
	response, err := client.requestApiAndErrorRetry(url, "GET", nil)
	if err != nil {
		return nil, fmt.Errorf("接口请求失败, err: %s", err)
	}

	var result define.PluginLogResult

	err = json.Unmarshal(response, &result)

	if err != nil {
		return nil, fmt.Errorf("json 反序列化失败, err: %s", err)
	}
	return &result, nil
}

func (client *Client) GetPluginList(limit int, offset int) (*define.PluginListReturn, error) {
	// 获取插件列表
	if !client.UsePluginService {
		return nil, fmt.Errorf("插件服务未开启")
	}

	pathParams := fmt.Sprintf("system/bk_plugins?limit=%d&offset=%d", limit, offset)

	url := client.preparePassApiRequest(pathParams)

	response, err := client.requestApiAndErrorRetry(url, "GET", nil)
	if err != nil {
		return nil, fmt.Errorf("url返回参数读取失败, err:%s", err)
	}

	var result define.PluginListReturn
	err = json.Unmarshal(response, &result)

	if err != nil {
		return nil, fmt.Errorf("url返回参数读取失败, err:%s", err)
	}
	return &result, nil
}
