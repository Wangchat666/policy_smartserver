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

func (s *StateCouncilContentColly) getRules() []*service.Rule {
	return []*service.Rule{
		s.zhengceCollector(),
	}
}

func (s *StateCouncilContentColly) updateTitle(e *colly.HTMLElement) {
	title := utils.TidyString(e.Text)
	s.metaDal.UpdateMetaTitle(title, e.Request.URL.String())
}

func (s *StateCouncilContentColly) updateContent(e *colly.HTMLElement) {
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
	}
	if meta != nil {
		sdIDs := s.dMapDal.GetDepartmentIDsByMetaID(meta.ID)
		es.IndexDoc(meta.Date, meta.DepartmentID, meta.ProvinceID, meta.Title, meta.Url, string(text), sdIDs)
		fmt.Println("education content updated") //我添加的检查是否上传es
	} else {
		fmt.Println("meta未查询到！！")
	}

}

func (s *StateCouncilContentColly) zhengceCollector() *service.Rule {

	rule := regexp.MustCompile("https?://www\\.gov\\.cn/zhengce/.*\\.html?")

	hfContent := &service.HtmlFunc{
		QuerySelect: "#UCAP-CONTENT",
		F:           s.updateContent,
	}

	return service.NormalRule(rule, hfContent)
}
