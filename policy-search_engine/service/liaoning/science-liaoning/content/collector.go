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
		s.kjtCollector(),
		s.weixinCollector(),
		s.sjzuCollector(),
		s.wapCollector(),
		s.newsCollector(),
		s.liaoningCollector(),
		s.newslndCollector(),
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

func (s *ScienceContentColly) kjtCollector() *service.Rule {
	//https://kjt.ln.gov.cn/kjt/kjzc/lnkjzc/2023121515172567535/index.shtml
	//https://mp.weixin.qq.com/s/C39teFFzSjKgPrllZD8MyA
	//https://kjt.ln.gov.cn/uiFramework/js/pdfjs/web/viewer.html?file=/eportal/fileDir/data/lnskjt/P020211230678959181796.pdf

	//地方政策
	//https://kjt.ln.gov.cn/kjt/kjzc/dfkjzc/ABB393B2E3A94DAA99091034F424C63F/index.shtml
	//https://kjt.ln.gov.cn/uiFramework/js/pdfjs/web/viewer.html?file=/eportal/fileDir/data/lnskjt/P020191010695729830609.pdf
	rule := regexp.MustCompile("https?://kjt\\.ln\\.gov\\.cn/kjt/.*\\.shtml?")

	hfTitle := &service.HtmlFunc{
		QuerySelect: ".govXLTitle",
		F:           s.updateTitle,
	}

	hfContent := &service.HtmlFunc{
		QuerySelect: ".TRS_Editor",
		F:           s.updateContent,
	}

	return service.NormalRule(rule, hfTitle, hfContent)
}

func (s *ScienceContentColly) weixinCollector() *service.Rule {
	//https://mp.weixin.qq.com/s/x7zieI_HtcpRxuBb4vRC-g
	rule := regexp.MustCompile("https?://mp\\.weixin\\.qq\\.com/.*")

	hfTitle := &service.HtmlFunc{
		QuerySelect: ".rich_media_title ",
		F:           s.updateTitle,
	}

	hfContent := &service.HtmlFunc{
		QuerySelect: "#js_content",
		F:           s.updateContent,
	}

	return service.NormalRule(rule, hfTitle, hfContent)
}
func (s *ScienceContentColly) sjzuCollector() *service.Rule {
	//https://www.sjzu.edu.cn/info/1381/84231.htm
	rule := regexp.MustCompile("https?://www\\.sjzu\\.edu\\.cn/.*\\.htm?")

	hfTitle := &service.HtmlFunc{
		QuerySelect: ".xqnr_tit",
		F:           s.updateTitle,
	}

	hfContent := &service.HtmlFunc{
		QuerySelect: ".v_news_content",
		F:           s.updateContent,
	}

	return service.NormalRule(rule, hfTitle, hfContent)
}
func (s *ScienceContentColly) wapCollector() *service.Rule {
	//https://wap.lnrbxmt.com/video_details.html?from=androidapp&id=369314&timestamp=98595032647524426
	rule := regexp.MustCompile("https?://wap\\.lnrbxmt\\.com/.*")

	hfTitle := &service.HtmlFunc{
		QuerySelect: ".article-top-inner",
		F:           s.updateTitle,
	}

	hfContent := &service.HtmlFunc{
		QuerySelect: ".details",
		F:           s.updateContent,
	}

	return service.NormalRule(rule, hfTitle, hfContent)
}
func (s *ScienceContentColly) newsCollector() *service.Rule {
	//http://www.news.cn/politics/leaders/2022-12/24/c_1129230368.htm
	rule := regexp.MustCompile("https?://www\\.news\\.cn/.*")

	hfTitle := &service.HtmlFunc{
		QuerySelect: ".title",
		F:           s.updateTitle,
	}

	hfContent := &service.HtmlFunc{
		QuerySelect: "#detail",
		F:           s.updateContent,
	}

	return service.NormalRule(rule, hfTitle, hfContent)
}
func (s *ScienceContentColly) liaoningCollector() *service.Rule {
	//http://liaoning.nen.com.cn/network/liaoningnews/lnnewsyuanchuang/2023/12/05/586651035934660109.shtml
	rule := regexp.MustCompile("https?://liaoning\\.nen\\.com\\.cn/.*")

	hfTitle := &service.HtmlFunc{
		QuerySelect: ".xwzw_title",
		F:           s.updateTitle,
	}

	hfContent := &service.HtmlFunc{
		QuerySelect: ".xwzw_t2",
		F:           s.updateContent,
	}

	return service.NormalRule(rule, hfTitle, hfContent)
}
func (s *ScienceContentColly) newslndCollector() *service.Rule {
	//https://news.lnd.com.cn/system/2023/09/08/030432462.shtml
	rule := regexp.MustCompile("https?://www\\.news\\.lnd\\.com\\.cn/.*")

	hfTitle := &service.HtmlFunc{
		QuerySelect: ".yahei newstittle",
		F:           s.updateTitle,
	}

	hfContent := &service.HtmlFunc{
		QuerySelect: ".news",
		F:           s.updateContent,
	}

	return service.NormalRule(rule, hfTitle, hfContent)
}
