# BK-PLUGIN-CLIENT-GO

蓝鲸插件 go 语言对接SDK。通过该SDK，你可以实现:
- 获取插件列表
- 获取插件详情
- 执行插件
- 查询插件状态
- 获取插件日志

## 环境变量
```text
USE_PLUGIN_SERVICE = 0 or 1 是否开启sdk 服务
PLUGIN_SERVICE_APIGW_APP_CODE = 调用身份的app_code
PLUGIN_SERVICE_APIGW_APP_SECRET = 调用身份的 app_seret
APIGW_ENVIRONMENT = 环境, prod or stage
APIGW_NETWORK_PROTOCAL = 协议 https or http
APIGW_URL_SUFFIX = APIGW的前缀地址
BKAPP_INVOKE_PAAS_RETRY_NUM = 失败重试次数
```

## 使用样例
```go
client := api.GetClient("plugin-demo")
// 查询插件列表
results, _ := client.GetPluginList(1, 2)
// 获取插件app详情, client初始化需要传入 plugin_code 才能获取详情
client.GetPluginAppDetail()
// 获取插件元数据
client.GetPluginMeta()
// 执行插件
client.Invoke(version, data)
// 获取插件状态
client.GetSchedule(trace_id)
// 获取插件日志
GetPluginLogs(traceId string, scrollId interface{})
```
