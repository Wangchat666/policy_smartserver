package content

import (
	"PolicySearchEngine/dao/es"
	"PolicySearchEngine/service"
	"PolicySearchEngine/utils"
	"fmt"
	"github.com/gocolly/colly"
	"regexp"
	"strings"
)

func (s *ScienceContentColly) getRules() []*service.Rule {
	return []*service.Rule{
		s.kjtCollector(),
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
	rule := regexp.MustCompile("https?://kjt\\.ln\\.gov\\.cn/kjt/kjzc/dfkjzc/.*\\.shtml?")

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

func (s *ScienceContentColly) xxgkCollector() *service.Rule {
	rule := regexp.MustCompile("https?://www\\.most\\.gov\\.cn/xxgk/.*\\.html?")

	hfTitle := &service.HtmlFunc{
		QuerySelect: ".xxgk_title",
		F:           s.updateTitle,
	}

	hfContent := &service.HtmlFunc{
		QuerySelect: "#Zoom",
		F:           s.updateContent,
	}

	return service.NormalRule(rule, hfTitle, hfContent)
}
