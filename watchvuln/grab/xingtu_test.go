package grab

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestXingtu(t *testing.T) {
	assert := require.New(t)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*300)
	defer cancel()

	grab := NewXingtuCrawler()
	vulns, err := grab.GetUpdate(ctx, 2) // 测试获取2页数据
	assert.Nil(err)

	count := 0
	for _, v := range vulns {
		t.Logf("get vuln info %s", v)
		count++
		// 基本字段验证
		assert.NotEmpty(v.UniqueKey)
		assert.NotEmpty(v.Title)
		assert.NotEmpty(v.Description)
		assert.NotEmpty(v.Disclosure)
		assert.NotEmpty(v.From)

		// 验证 UniqueKey 格式
		assert.Contains(v.UniqueKey, "DAS-T")

		// 验证 From 字段格式
		assert.Contains(v.From, "https://ti.dbappsecurity.com.cn/vul/")

		// 验证标签非空（如果有影响信息）
		if len(v.Tags) > 0 {
			for _, tag := range v.Tags {
				assert.NotEmpty(tag)
			}
		}
	}

	// 确保获取到漏洞信息
	assert.Greater(count, 0)

	// 验证创建者信息
	provider := grab.ProviderInfo()
	assert.Equal("xingtu", provider.Name)
	assert.Equal("安恒安全星图平台", provider.DisplayName)
	assert.Equal("https://ti.dbappsecurity.com.cn/vul", provider.Link)
}

// 测试严重等级过滤
func TestXingtuValuable(t *testing.T) {
	assert := require.New(t)
	grab := NewXingtuCrawler()

	// 测试各种严重等级
	testCases := []struct {
		severity Severity
		valuable bool
	}{
		{Critical, true},
		{High, true},
		{Medium, false},
		{Low, false},
	}

	for _, tc := range testCases {
		vuln := &VulnInfo{Severity: tc.severity}
		assert.Equal(tc.valuable, grab.IsValuable(vuln))
	}
}
