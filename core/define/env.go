package define

import (
	"os"
)

func getEnv(key string, defalutValue string) string {

	val := os.Getenv(key)
	if len(val) != 0 {
		return val
	}
	return defalutValue

}

func getEnvReturnBool(key string, defalutValue bool) bool {
	val := os.Getenv(key)
	if len(val) != 0 {
		return val == "1"
	}
	return defalutValue

}

var UsePluginService = getEnvReturnBool("USE_PLUGIN_SERVICE", false)            // 是否开启插件服务
var PluginServiceApiGwAppCode = getEnv("PLUGIN_SERVICE_APIGW_APP_CODE", "")     // APP_CODE
var PluginServiceApiGwAppSecret = getEnv("PLUGIN_SERVICE_APIGW_APP_SECRET", "") // APP_SECRET
var ApiGwEnvironment = getEnv("APIGW_ENVIRONMENT", "")
var ApiGwNetworkProtocal = getEnv("APIGW_NETWORK_PROTOCAL", "http")
var ApiGwUrlSuffix = getEnv("APIGW_URL_SUFFIX", "")
var BkAppInvokePaasRetryNum = getEnv("BKAPP_INVOKE_PAAS_RETRY_NUM", "3")
var ApiGwUserAuthKeyName = getEnv("BKAPP_INVOKE_PAAS_RETRY_NUM", "bk_token")

var DefaultHostType = 2
