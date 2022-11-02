package api

import (
	"testing"
)

func TestGetClient(t *testing.T) {

	client, _ := GetClient("")
	if client.UsePluginService {
		t.Errorf("client 初始化失败，UsePluginService 不符合预期")
	}
}
