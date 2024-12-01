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

func (s *EducationContentColly) getRules() []*service.Rule {
	return []*service.Rule{
		s.zcfgCollector(),
		s.srcsiteCollector(),
	}
}

func (s *EducationContentColly) updateTitle(e *colly.HTMLElement) {
	title := utils.TidyString(e.Text)
	s.metaDal.UpdateMetaTitle(title, e.Request.URL.String())
}

func (s *EducationContentColly) updateContent(e *colly.HTMLElement) {
	var text []byte
	e.ForEach("*", func(_ int, child *colly.HTMLElement) {
		label := strings.ToLower(child.Name)
		if label == "style" || label == "table" || label == "script" {
			return
		}
		text = append(text, []byte(child.Text)...)
	})
	s.contentDal.InsertContent(e.Request.URL.String(), string(text)) //将提取的文本内容插入到数据库中，使用当前请求的 URL 作为标识。

	meta := s.metaDal.GetMetaByUrl(e.Request.URL.String()) //尝试从数据库中获取与当前 URL 对应的元数据
	if meta == nil {
		meta = s.metaDal.GetMetaByUrl(e.Request.Headers.Get("Referer"))
	}
	if meta != nil {
		sdIDs := s.dMapDal.GetDepartmentIDsByMetaID(meta.ID)
		es.IndexDoc(meta.Date, meta.DepartmentID, meta.ProvinceID, meta.Title, meta.Url, string(text), sdIDs)
		fmt.Println("education content updated") //我添加的检查是否上传es
		return
	}
}

func (s *EducationContentColly) zcfgCollector() *service.Rule {

	rule := regexp.MustCompile("https?://www\\.moe\\.gov\\.cn/jyb_sjzl/sjzl_zcfg/.*\\.html?")

	hfTitle := &service.HtmlFunc{
		QuerySelect: ".moe-detail-box h1",
		F:           s.updateTitle,
	}

	hfContent := &service.HtmlFunc{
		QuerySelect: ".TRS_Editor",
		F:           s.updateContent,
	}

	return service.NormalRule(rule, hfTitle, hfContent)
}

func (s *EducationContentColly) srcsiteCollector() *service.Rule {

	rule := regexp.MustCompile("https?://www\\.moe\\.gov\\.cn/srcsite/.*\\.html?")

	hfTitle := &service.HtmlFunc{
		QuerySelect: ".details-policy-box h1",
		F:           s.updateTitle,
	}

	hfContent := &service.HtmlFunc{
		QuerySelect: ".details-policy-box",
		F:           s.updateContent,
	}

	return service.NormalRule(rule, hfTitle, hfContent)
}
