package grab

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestShuZiGuanXing(t *testing.T) {
	assert := require.New(t)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*300)
	defer cancel()

	grab := NewShuZiGuanXingCrawler()
	vulns, err := grab.GetUpdate(ctx, 5)
	assert.Nil(err)

	for _, v := range vulns {
		t.Logf("get vuln info %s", v)
		assert.NotEmpty(v.UniqueKey)
		assert.NotEmpty(v.Title)
		assert.NotEmpty(v.Disclosure)
		assert.NotEmpty(v.From)
	}
}
