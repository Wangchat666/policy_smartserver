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

func (s *IndustryInformatizationContentColly) getRules() []*service.Rule {
	return []*service.Rule{
		s.miitCollector(),
	}
}

func (s *IndustryInformatizationContentColly) updateTitle(e *colly.HTMLElement) {
	title := utils.TidyString(e.Text)
	s.metaDal.UpdateMetaTitle(title, e.Request.URL.String())
}

func (s *IndustryInformatizationContentColly) updateContent(e *colly.HTMLElement) {
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
		fmt.Println("未上传") //我添加的检查是否上传es
	}
	if meta != nil {
		sdIDs := s.dMapDal.GetDepartmentIDsByMetaID(meta.ID)
		es.IndexDoc(meta.Date, meta.DepartmentID, meta.ProvinceID, meta.Title, meta.Url, string(text), sdIDs)
		fmt.Println("education content updated") //我添加的检查是否上传es
	}
}

func (s *IndustryInformatizationContentColly) miitCollector() *service.Rule {

	rule := regexp.MustCompile("https?://wap\\.miit\\.gov\\.cn/.*\\.html?")

	hfTitle := &service.HtmlFunc{
		QuerySelect: "#con_title",
		F:           s.updateTitle,
	}

	hfContent := &service.HtmlFunc{
		QuerySelect: "#con_con",
		F:           s.updateContent,
	}

	return service.NormalRule(rule, hfTitle, hfContent)
}
