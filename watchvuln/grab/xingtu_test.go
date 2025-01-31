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
	vulns, err := grab.GetUpdate(ctx, 5)
	assert.Nil(err)

	for _, v := range vulns {
		t.Logf("get vuln info %s", v)
		assert.NotEmpty(v.UniqueKey)
		assert.NotEmpty(v.Title)
		assert.NotEmpty(v.Description)
		assert.NotEmpty(v.From)
		// 打印补丁修复信息以验证是否正确解析
		t.Logf("solutions: %s", v.Solutions)
	}
}
