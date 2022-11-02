package define

type Plugin struct {
	// 插件基本信息
	Code    string      `json:"code,omitempty"`
	Name    string      `json:"name,omitempty"`
	Id      string      `json:"id,omitempty"`
	LogoUrl string      `json:"logo_url,omitempty"`
	Creator string      `json:"creator,omitempty"`
	TagInfo interface{} `json:"tag_info,omitempty"`
	Created string      `json:"created,omitempty"`
	Region  string      `json:"region,omitempty"`
}

type PluginListReturn struct {
	// 插件列表返回
	Count   int
	Results []Plugin
}

type PluginAppDetail struct {
	// 插件详情
	Url       string
	Urls      []string
	Name      string
	Code      string
	ApiGwName string
}

type PluginMeta struct {
	// 插件元信息
	Code             string
	Versions         []string
	Language         string
	Description      string
	FrameworkVersion string
	RuntimeVersion   string
}

type PluginMetaReturn struct {
	// 插件元信息返回
	Result  bool
	Data    PluginMeta
	Message string
}

type PluginDetail struct {
	// 插件详情
	Desc          string
	Version       string
	Inputs        map[string]interface{}
	OutPuts       map[string]interface{}
	ContextInputs map[string]interface{}
}

type PluginDetailReturn struct {
	// 插件详情返回
	Result  bool
	Data    PluginDetail
	Message string
}

type PluginInvokeResultData struct {
	Outputs map[string]interface{}
	State   int
	Err     string
}

type PluginInvokeResult struct {
	Result  bool                   `json:"result,omitempty"`
	Data    PluginInvokeResultData `json:"data"`
	Message string                 `json:"message,omitempty"`
	TraceId string                 `json:"trace_id,omitempty"`
}

type PluginScheduleResultData struct {
	TraceId       string `json:"trace_id,omitempty"`
	PluginVersion string `json:"plugin_version,omitempty"`
	State         int    `json:"state,omitempty"`
	Err           string `json:"err,omitempty"`
	CreateAt      string `json:"create_at,omitempty"`
	FinishAt      string `json:"finish_at,omitempty"`
}

type PluginScheduleResult struct {
	Result  bool                     `json:"result,omitempty"`
	Data    PluginScheduleResultData `json:"data"`
	Message string                   `json:"message,omitempty"`
	TraceId string                   `json:"trace_id,omitempty"`
}

type PluginLogData struct {
	PluginCode  string                 `json:"plugin_code,omitempty"`
	Environment string                 `json:"environment,omitempty"`
	ProcessId   string                 `json:"process_id,omitempty"`
	Stream      string                 `json:"stream,omitempty"`
	Message     string                 `json:"message,omitempty"`
	Detail      map[string]interface{} `json:"detail,omitempty"`
	Ts          string                 `json:"ts,omitempty"`
}

type PluginLogResult struct {
	ScrollId string          `json:"scroll_id,omitempty"`
	Logs     []PluginLogData `json:"logs,omitempty"`
	Total    int             `json:"total,omitempty"`
}
