package grab

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/imroc/req/v3"
	"github.com/kataras/golog"
	"github.com/zema1/watchvuln/util"
)

type XingtuCrawler struct {
	client *req.Client
	log    *golog.Logger
}

func NewXingtuCrawler() Grabber {
	client := util.NewHttpClient()
	client.SetCommonHeader("Referer", "https://ti.dbappsecurity.com.cn/")

	return &XingtuCrawler{
		client: client,
		log:    golog.Child("[xingtu]"),
	}
}

func (x *XingtuCrawler) ProviderInfo() *Provider {
	return &Provider{
		Name:        "xingtu",
		DisplayName: "安恒安全星图平台",
		Link:        "https://ti.dbappsecurity.com.cn/vul",
	}
}

func (x *XingtuCrawler) GetUpdate(ctx context.Context, pageLimit int) ([]*VulnInfo, error) {
	if pageLimit > 4 {
		pageLimit = 4
	}

	var results []*VulnInfo
	// 编译正则表达式用于提取DAS-T编号
	dasPattern := regexp.MustCompile(`DAS-T\d{6}`)

	for i := 1; i <= pageLimit; i++ {
		select {
		case <-ctx.Done():
			return results, ctx.Err()
		default:
		}

		// 获取列表页面
		url := fmt.Sprintf("https://ti.dbappsecurity.com.cn/vul?page=%d", i)
		resp, err := x.client.R().SetContext(ctx).Get(url)
		if err != nil {
			x.log.Errorf("Failed to get page %d: %v", i, err)
			continue
		}

		// 从页面内容中提取所有DAS-T编号
		content := resp.String()
		matches := dasPattern.FindAllString(content, -1)

		// 对每个找到的DAS编号获取详细信息
		pageVulns := 0
		for _, dasID := range matches {
			// 获取漏洞详情
			vulnInfo, err := x.parseVulnDetail(ctx, dasID)
			if err != nil {
				x.log.Warnf("Failed to parse vuln detail for %s: %v", dasID, err)
				continue
			}

			if vulnInfo != nil {
				results = append(results, vulnInfo)
				pageVulns++
			}
		}

		x.log.Infof("got %d vulns from page %d", pageVulns, i)
	}

	return results, nil
}

func (x *XingtuCrawler) parseVulnDetail(ctx context.Context, dasID string) (*VulnInfo, error) {
	url := fmt.Sprintf("https://ti.dbappsecurity.com.cn/vul/%s", dasID)
	resp, err := x.client.R().SetContext(ctx).Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get detail page: %v", err)
	}

	content := resp.String()

	// 正则表达式
	titleRe := regexp.MustCompile(`<div class="hole-name"[^>]*>(.*?)</div>`)
	levelRe := regexp.MustCompile(`<div class="level"[^>]*>(.*?)</div>`)
	descRe := regexp.MustCompile(`<div class="detail-desc"[^>]*>(.*?)</div>`)
	cveRe := regexp.MustCompile(`CVE编号：(CVE-\d{4}-\d+)`)
	timeRe := regexp.MustCompile(`发布时间：(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2})`)
	vulnResultRe := regexp.MustCompile(`<span class="hole-result"[^>]*>([^<]+)</span>`)
	solutionRe := regexp.MustCompile(`(?s)<div class="detail-box detail-box-6".*?<div class="detail-desc"[^>]*>(.*?)</div>`)

	// 提取信息
	title := ""
	if matches := titleRe.FindStringSubmatch(content); len(matches) > 1 {
		title = strings.TrimSpace(matches[1])
	}

	// 提取严重级别
	severity := Low
	if matches := levelRe.FindStringSubmatch(content); len(matches) > 1 {
		switch strings.TrimSpace(matches[1]) {
		case "超危":
			severity = Critical
		case "高危":
			severity = High
		case "中危":
			severity = Medium
		case "低危":
			severity = Low
		}
	}

	// 提取描述
	description := ""
	if matches := descRe.FindStringSubmatch(content); len(matches) > 1 {
		description = strings.TrimSpace(matches[1])
	}

	// 提取CVE编号
	cveID := ""
	if matches := cveRe.FindStringSubmatch(content); len(matches) > 1 {
		cveID = matches[1]
	}

	// 提取发布时间
	disclosure := ""
	if matches := timeRe.FindStringSubmatch(content); len(matches) > 1 {
		if t, err := time.Parse("2006-01-02 15:04:05", matches[1]); err == nil {
			disclosure = t.Format("2006-01-02")
		}
	}

	// 提取标签（只包含漏洞影响）
	var tags []string
	vulnResultMatches := vulnResultRe.FindAllStringSubmatch(content, -1)
	for _, match := range vulnResultMatches {
		if len(match) > 1 {
			result := strings.TrimSpace(match[1])
			if result != "" {
				tags = append(tags, result)
			}
		}
	}

	// 提取补丁修复信息
	solutions := ""
	if matches := solutionRe.FindStringSubmatch(content); len(matches) > 1 {
		solutions = strings.TrimSpace(matches[1])
	}

	// 如果没有找到标题，说明页面可能无效
	if title == "" {
		return nil, fmt.Errorf("empty title, page might be invalid")
	}

	vulnInfo := &VulnInfo{
		UniqueKey:   dasID,
		Title:       title,
		Description: description,
		Severity:    severity,
		CVE:         cveID,
		Disclosure:  disclosure,
		Solutions:   solutions,
		From:        url,
		Tags:        tags,
		Creator:     x,
	}

	x.log.Debugf("Successfully parsed vuln: %s - %s, tags: %v", dasID, title, tags)
	return vulnInfo, nil
}

func (x *XingtuCrawler) IsValuable(info *VulnInfo) bool {
	return info.Severity == High || info.Severity == Critical
}
