package content

import (
	"fmt"
	"github.com/gocolly/colly"
	"policy-search_engine/dao/es"
	"policy-search_engine/service"
	"policy-search_engine/utils"
	"regexp"
	"strings"
)

func (s *ScienceContentColly) getRules() []*service.Rule {
	return []*service.Rule{
		s.TZGG2079Collector(),
		s.ZXGZ7816Collector(),
		s.JZBGHYT7299Collector(),
		s.KJGH20201119Collector(),
		s.XZZF201019Collector(),
		s.ZFXXGK2908_1Collector(),
		s.managecolCollector(),
		s.ZCFG148_1Collector(),
	}
}

func (s *ScienceContentColly) updateTitle(e *colly.HTMLElement) {
	title := utils.TidyString(e.Text)
	s.metaDal.UpdateMetaTitle(title, e.Request.URL.String())
}

func (s *ScienceContentColly) updateContent(e *colly.HTMLElement) {
	var text []byte
	e.ForEach("*", func(_ int, child *colly.HTMLElement) {
		label := strings.ToLower(child.Name)
		if label == "style" || label == "table" || label == "script" {
			return
		}
		text = append(text, []byte(child.Text)...)
	})
	s.contentDal.InsertContent(e.Request.URL.String(), string(text))

	meta := s.metaDal.GetMetaByUrl(e.Request.URL.String())
	if meta == nil {
		meta = s.metaDal.GetMetaByUrl(e.Request.Headers.Get("Referer"))
		fmt.Println("未上传")
	}
	if meta != nil {
		sdIDs := s.dMapDal.GetDepartmentIDsByMetaID(meta.ID)
		es.IndexDoc(meta.Date, meta.DepartmentID, meta.ProvinceID, meta.Title, meta.Url, string(text), sdIDs)
		fmt.Println("education content updated") //我添加的检查是否上传es

	}
}

func (s *ScienceContentColly) TZGG2079Collector() *service.Rule {
	//https://kxjs.tj.gov.cn/ZWGK4143/TZGG2079/202407/t20240711_6673541.html

	rule := regexp.MustCompile("https?://kxjs\\.tj\\.gov\\.cn/ZWGK4143/TZGG2079/\\d{6}/t\\d{8}_\\d+.html")

	hfTitle := &service.HtmlFunc{
		QuerySelect: ".common-content-mainTitle",
		F:           s.updateTitle,
	}

	hfContent := &service.HtmlFunc{
		QuerySelect: ".page_info",
		F:           s.updateContent,
	}

	return service.NormalRule(rule, hfTitle, hfContent)
}

func (s *ScienceContentColly) ZXGZ7816Collector() *service.Rule {
	//https://kxjs.tj.gov.cn/ZWGK4143/ZXGZ7816/KJJR3268/GZDT5932/202407/t20240701_6665348.html
	rule := regexp.MustCompile("https?://kxjs\\.tj\\.gov\\.cn/ZWGK4143/ZXGZ7816/KJJR3268/GZDT5932/\\d{6}/t\\d{8}_\\d+.html")

	hfTitle := &service.HtmlFunc{
		QuerySelect: ".common-content-mainTitle",
		F:           s.updateTitle,
	}

	hfContent := &service.HtmlFunc{
		QuerySelect: ".page_info",
		F:           s.updateContent,
	}

	return service.NormalRule(rule, hfTitle, hfContent)
}
func (s *ScienceContentColly) JZBGHYT7299Collector() *service.Rule {
	//https://kxjs.tj.gov.cn/ZWGK4143/JZBGHYT7229/202407/t20240708_6670276.html
	rule := regexp.MustCompile("https?://kxjs\\.tj\\.gov\\.cn/ZWGK4143/JZBGHYT7229/\\d{6}/t\\d{8}_\\d+.html")

	hfTitle := &service.HtmlFunc{
		QuerySelect: ".common-content-mainTitle",
		F:           s.updateTitle,
	}

	hfContent := &service.HtmlFunc{
		QuerySelect: ".page_info",
		F:           s.updateContent,
	}

	return service.NormalRule(rule, hfTitle, hfContent)
}
func (s *ScienceContentColly) KJGH20201119Collector() *service.Rule {
	//https://kxjs.tj.gov.cn/ZWGK4143/KJGH20201119/202404/t20240402_6589629.html
	rule := regexp.MustCompile("https?://kxjs\\.tj\\.gov\\.cn/ZWGK4143/KJGH20201119/\\d{6}/t\\d{8}_\\d+.html")

	hfTitle := &service.HtmlFunc{
		QuerySelect: ".common-content-mainTitle",
		F:           s.updateTitle,
	}

	hfContent := &service.HtmlFunc{
		QuerySelect: ".page_info",
		F:           s.updateContent,
	}

	return service.NormalRule(rule, hfTitle, hfContent)
}
func (s *ScienceContentColly) XZZF201019Collector() *service.Rule {
	//https://kxjs.tj.gov.cn/ZWGK4143/XZZF201019/fzzfjsndbg/202403/t20240327_6572155.html
	//https://kxjs.tj.gov.cn/ZWGK4143/XZZF201019/fzzfjsgz/202406/t20240603_6640964.html
	rule := regexp.MustCompile("https?://kxjs\\.tj\\.gov\\.cn/ZWGK4143/XZZF201019/\\D+/\\d{6}/t\\d{8}_\\d+.html")

	hfTitle := &service.HtmlFunc{
		QuerySelect: ".common-content-mainTitle",
		F:           s.updateTitle,
	}

	hfContent := &service.HtmlFunc{
		QuerySelect: ".page_info",
		F:           s.updateContent,
	}

	return service.NormalRule(rule, hfTitle, hfContent)
}
func (s *ScienceContentColly) ZFXXGK2908_1Collector() *service.Rule {
	//https://kxjs.tj.gov.cn/ZWGK4143/ZFXXGK2908_1/ZFXXGKNB5452/202012/t20201211_4860633.html
	rule := regexp.MustCompile("https?://kxjs\\.tj\\.gov\\.cn/ZWGK4143/ZFXXGK2908_1/ZFXXGKNB5452/\\d{6}/t\\d{8}_\\d+.html")

	hfTitle := &service.HtmlFunc{
		QuerySelect: ".qt-title",
		F:           s.updateTitle,
	}

	hfContent := &service.HtmlFunc{
		QuerySelect: ".view TRS_UEDITOR trs_paper_default trs_external",
		F:           s.updateContent,
	}

	return service.NormalRule(rule, hfTitle, hfContent)
}
func (s *ScienceContentColly) managecolCollector() *service.Rule {
	//https://kxjs.tj.gov.cn/managecol/ZCWJ0923/kjjzcwj09233/202012/t20201211_4858991.html
	rule := regexp.MustCompile("https?://kxjs\\.tj\\.gov\\.cn/managecol/ZCWJ0923/kjjzcwj09233/\\d{6}/t\\d{8}_\\d+.html")

	hfTitle := &service.HtmlFunc{
		QuerySelect: ".sx-con",
		F:           s.updateTitle,
	}

	hfContent := &service.HtmlFunc{
		QuerySelect: ".xl-zw-cons",
		F:           s.updateContent,
	}

	return service.NormalRule(rule, hfTitle, hfContent)
}
func (s *ScienceContentColly) ZCFG148_1Collector() *service.Rule {
	//https://kxjs.tj.gov.cn/ZWGK4143/ZCFG148_1/ZCFB4222/2015N2998/202012/t20201211_4860723.html
	rule := regexp.MustCompile("https?://kxjs\\.tj\\.gov\\.cn/ZWGK4143/ZCFG148_1/ZCFB4222/2015N2998/\\d{6}/t\\d{8}_\\d+.html")

	hfTitle := &service.HtmlFunc{
		QuerySelect: ".common-content-mainTitle",
		F:           s.updateTitle,
	}

	hfContent := &service.HtmlFunc{
		QuerySelect: ".page_info",
		F:           s.updateContent,
	}

	return service.NormalRule(rule, hfTitle, hfContent)
}
