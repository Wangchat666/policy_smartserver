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
		//s.webCollector(),
		s.shenyangCollector(),
		//s.dlCollector(),
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
		fmt.Println(" content updated") //我添加的检查是否上传es
	}
}

func (s *ScienceContentColly) webCollector() *service.Rule {
	//https://www.ln.gov.cn/web/zwgkx/zfwj/szfwj/zfwj2011_125195/C8C3B9765A004FEE800CB7F391EF9388/index.shtml
	//https://kjt.ln.gov.cn/kjt/kjgz/gxkj/C2EF3FD59E6D4E18A6E9499B491F07CE/index.shtml
	rule := regexp.MustCompile("https?://www\\.ln\\.gov\\.cn/web/.*\\.shtml?")

	hfTitle := &service.HtmlFunc{
		QuerySelect: ".title",
		F:           s.updateTitle,
	}

	hfContent := &service.HtmlFunc{
		QuerySelect: ".TRS_Editor",
		F:           s.updateContent,
	}

	return service.NormalRule(rule, hfTitle, hfContent)
}
func (s *ScienceContentColly) shenyangCollector() *service.Rule {
	//https://www.shenyang.gov.cn/zwgk/zcwj/fggz/sdfxfgzl/202308/t20230828_4517478.html
	rule := regexp.MustCompile("https?://www\\.shenyang\\.gov\\.cn/.*\\.html?")

	hfTitle := &service.HtmlFunc{
		QuerySelect: ".dlist_title",
		F:           s.updateTitle,
	}

	hfContent := &service.HtmlFunc{
		QuerySelect: ".view TRS_UEDITOR trs_paper_default trs_external trs_web",
		F:           s.updateContent,
	}

	return service.NormalRule(rule, hfTitle, hfContent)
}
func (s *ScienceContentColly) dlCollector() *service.Rule {
	//https://www.dl.gov.cn/art/2022/10/27/art_852_2050004.html?xxgkhide=1
	rule := regexp.MustCompile("https?://www\\.dl\\.gov\\.cn/.*")

	hfTitle := &service.HtmlFunc{
		QuerySelect: ".text-center line-height-2",
		F:           s.updateTitle,
	}

	hfContent := &service.HtmlFunc{
		QuerySelect: "#zoom",
		F:           s.updateContent,
	}

	return service.NormalRule(rule, hfTitle, hfContent)
}
