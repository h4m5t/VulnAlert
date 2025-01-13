package grab

import (
	"context"
	"fmt"
	"strings"

	"github.com/imroc/req/v3"
	"github.com/kataras/golog"
	"github.com/zema1/watchvuln/util"
)

type ShuZiGuanXingCrawler struct {
	client *req.Client
	log    *golog.Logger
}

func NewShuZiGuanXingCrawler() Grabber {
	client := util.WrapApiClient(util.NewHttpClient())
	client.SetCommonHeader("Content-Type", "application/json")

	c := &ShuZiGuanXingCrawler{
		client: client,
		log:    golog.Child("[shuziguanxing]"),
	}
	return c
}

func (t *ShuZiGuanXingCrawler) ProviderInfo() *Provider {
	return &Provider{
		Name:        "shuziguanxing",
		DisplayName: "数字观星漏洞库",
		Link:        "https://poc.shuziguanxing.com",
	}
}

func (t *ShuZiGuanXingCrawler) GetUpdate(ctx context.Context, pageLimit int) ([]*VulnInfo, error) {
	// 强制设置为3页，无视传入的pageLimit参数
	pageLimit = 3 //只爬取前3页

	var results []*VulnInfo

	// 获取高危(3)和严重(4)级别的漏洞
	for _, severity := range []int{3, 4} {
		for page := 1; page <= pageLimit; page++ {
			select {
			case <-ctx.Done():
				return results, ctx.Err()
			default:
			}

			pageResult, err := t.getPage(ctx, page, severity)
			if err != nil {
				return results, err
			}
			t.log.Infof("got %d vulns from severity %d page %d", len(pageResult), severity, page)
			results = append(results, pageResult...)
		}
	}

	return results, nil
}

func (t *ShuZiGuanXingCrawler) getPage(ctx context.Context, pageNum int, severity int) ([]*VulnInfo, error) {
	payload := map[string]interface{}{
		"issueExtWhere": map[string]interface{}{
			"severity":       severity,
			"peopleSortType": nil,
			"timeSortType":   2,
			"type":           nil,
		},
		"pageBasic": map[string]interface{}{
			"numPerPage": 10,
			"pageNum":    pageNum,
		},
	}

	var resp struct {
		State string `json:"state"`
		Info  struct {
			Result []struct {
				ID       int    `json:"id"`
				AddTime  string `json:"addTime"`
				Name     string `json:"name"`
				Severity int    `json:"severity"`
				CveId    string `json:"cveId"`
			} `json:"result"`
		} `json:"info"`
	}

	_, err := t.client.R().
		SetContext(ctx).
		SetBody(payload).
		SetSuccessResult(&resp).
		Post("https://poc.shuziguanxing.com/pocweb/issueWarehouse/list")

	if err != nil {
		return nil, err
	}

	if resp.State != "success" {
		return nil, fmt.Errorf("api return error: %s", resp.State)
	}

	var results []*VulnInfo
	for _, item := range resp.Info.Result {
		// 获取详情
		detail, err := t.getVulnDetail(ctx, item.ID)
		if err != nil {
			t.log.Warnf("failed to get detail for %d: %s", item.ID, err)
			continue
		}

		// 转换严重级别
		severity := Low
		switch item.Severity {
		case 4:
			severity = Critical
		case 3:
			severity = High
		}

		var tags []string
		if detail.Source != "" {
			tags = append(tags, detail.Source)
		}

		info := &VulnInfo{
			UniqueKey:   fmt.Sprintf("DSO-%d", item.ID),
			Title:       item.Name,
			Description: stripHTML(detail.Info),
			Severity:    severity,
			CVE:         item.CveId,
			Disclosure:  detail.FindTime,
			Solutions: fmt.Sprintf("临时解决方案:\n%s\n\n官方解决方案:\n%s",
				stripHTML(detail.TemporarySolve),
				stripHTML(detail.OfficialSolve)),
			References: []string{detail.ReferenceUrl},
			Tags:       tags,
			From:       fmt.Sprintf("https://poc.shuziguanxing.com/pocweb/issueWarehouse/get?id=%d", item.ID),
			Creator:    t,
		}
		results = append(results, info)
	}

	return results, nil
}

func (t *ShuZiGuanXingCrawler) getVulnDetail(ctx context.Context, id int) (*vulnDetail, error) {
	var resp struct {
		State string     `json:"state"`
		Info  vulnDetail `json:"info"`
	}

	_, err := t.client.R().
		SetContext(ctx).
		SetSuccessResult(&resp).
		Get(fmt.Sprintf("https://poc.shuziguanxing.com/pocweb/issueWarehouse/get?id=%d", id))

	if err != nil {
		return nil, err
	}

	if resp.State != "success" {
		return nil, fmt.Errorf("api return error: %s", resp.State)
	}

	return &resp.Info, nil
}

type vulnDetail struct {
	Info           string `json:"info"`
	TemporarySolve string `json:"temporarySolve"`
	OfficialSolve  string `json:"officialSolve"`
	FindTime       string `json:"findTime"`
	Source         string `json:"source"`
	ReferenceUrl   string `json:"referenceUrl"`
}

func (t *ShuZiGuanXingCrawler) IsValuable(info *VulnInfo) bool {
	return info.Severity == Critical || info.Severity == High
}

// 简单的HTML标签清理
func stripHTML(s string) string {
	s = strings.ReplaceAll(s, "</p>", "\n")
	s = strings.ReplaceAll(s, "<br>", "\n")
	s = strings.ReplaceAll(s, "</br>", "\n")

	// 移除所有HTML标签
	for strings.Contains(s, "<") && strings.Contains(s, ">") {
		start := strings.Index(s, "<")
		end := strings.Index(s, ">")
		if start < end {
			s = s[:start] + s[end+1:]
		} else {
			break
		}
	}

	return strings.TrimSpace(s)
}
