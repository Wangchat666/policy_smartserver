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
		s.xxgkCollector(),
		s.kjzcCollector(),
		s.kjbgzCollector(),
		s.zhengceCollector(),
		s.gongbaoCollector(),
		s.xinwenCollector(),
		s.chinataxCollector(),
	}
}

func (s *ScienceContentColly) updateTitle(e *colly.HTMLElement) {
	title := utils.TidyString(e.Text)
	s.metaDal.UpdateMetaTitle(title, e.Request.URL.String())
}

func (s *ScienceContentColly) updateContent(e *colly.HTMLElement) {
	var text []byte
	e.ForEach("*", func(_ int, child *colly.HTMLElement) {
		label := strings.ToLower(child.Name) //将其标签名转换为小写
		if label == "style" || label == "table" || label == "script" {
			return
		}
		text = append(text, []byte(child.Text)...)
	})
	s.contentDal.InsertContent(e.Request.URL.String(), string(text)) //将提取的文本内容和当前请求的 URL 插入到数据库中

	meta := s.metaDal.GetMetaByUrl(e.Request.URL.String())
	if meta == nil {
		meta = s.metaDal.GetMetaByUrl(e.Request.Headers.Get("Referer"))
		fmt.Println("未上传")
	} //获取元数据，尝试通过当前请求的 URL 获取元数据。如果未找到元数据，则尝试通过请求头中的 Referer 字段获取元数据。
	if meta != nil {
		sdIDs := s.dMapDal.GetDepartmentIDsByMetaID(meta.ID)
		es.IndexDoc(meta.Date, meta.DepartmentID, meta.ProvinceID, meta.Title, meta.Url, string(text), sdIDs)
		fmt.Println("education content updated") //我添加的检查是否上传es
	} //如果找到了元数据，这段代码首先通过元数据的 ID 获取相关的部门 IDs，然后调用 es.IndexDoc 方法将元数据和提取的文本内容索引到 Elasticsearch 中。

} //主要功能是从 HTML 元素中提取文本内容，将其存储到数据库中，并更新相关的元数据和索引。

func (s *ScienceContentColly) xxgkCollector() *service.Rule {
	rule := regexp.MustCompile("https?://www\\.most\\.gov\\.cn/xxgk/.*\\.html?")

	hfTitle := &service.HtmlFunc{
		QuerySelect: ".xxgk_title", //QuerySelect 字段指定了一个 CSS 选择器
		F:           s.updateTitle, //F 字段指定了一个回调函数 s.updateTitle，用于处理提取到的标题。
	}

	hfContent := &service.HtmlFunc{
		QuerySelect: "#Zoom",
		F:           s.updateContent,
	}

	return service.NormalRule(rule, hfTitle, hfContent)
} //用于从特定的 URL 中提取标题和内容，并指定相应的处理函数。

func (s *ScienceContentColly) kjzcCollector() *service.Rule {

	rule1 := regexp.MustCompile("https?://www\\.most\\.gov\\.cn/satp/kjzc/zh/.*\\.html?")
	rule2 := regexp.MustCompile("https?://www\\.most\\.gov\\.cn/tztg/.*\\.html?")

	combinedRule := regexp.MustCompile(fmt.Sprintf(
		"(%s|%s)",
		rule1.String(),
		rule2.String(),
	))

	hfTitle := &service.HtmlFunc{
		QuerySelect: "#Title",
		F:           s.updateTitle,
	}

	hfContent := &service.HtmlFunc{
		QuerySelect: "#Zoom",
		F:           s.updateContent,
	}

	return service.NormalRule(combinedRule, hfTitle, hfContent)
}

func (s *ScienceContentColly) kjbgzCollector() *service.Rule {

	rule := regexp.MustCompile("https?://www\\.most\\.gov\\.cn/kjbgz/.*\\.html?")

	hfTitle := &service.HtmlFunc{
		QuerySelect: "#Title",
		F:           s.updateTitle,
	}

	hfContent := &service.HtmlFunc{
		QuerySelect: "#Zoom",
		F:           s.updateContent,
	}

	return service.NormalRule(rule, hfTitle, hfContent)
}

func (s *ScienceContentColly) zhengceCollector() *service.Rule {

	rule := regexp.MustCompile("https?://www\\.gov\\.cn/zhengce/content/.*\\.html?")

	hfTitle := &service.HtmlFunc{
		QuerySelect: "td[colspan='3']",
		F:           s.updateTitle,
	}

	hfContent := &service.HtmlFunc{
		QuerySelect: "#UCAP-CONTENT",
		F:           s.updateContent,
	}

	return service.NormalRule(rule, hfTitle, hfContent)
}

func (s *ScienceContentColly) gongbaoCollector() *service.Rule {

	rule := regexp.MustCompile("https?://www\\.gov\\.cn/gongbao/content/.*\\.html?")

	hfTitle := &service.HtmlFunc{
		QuerySelect: ".share-title",
		F:           s.updateTitle,
	}

	hfContent := &service.HtmlFunc{
		QuerySelect: ".pages_content",
		F:           s.updateContent,
	}

	return service.NormalRule(rule, hfTitle, hfContent)
}

func (s *ScienceContentColly) xinwenCollector() *service.Rule {

	rule := regexp.MustCompile("https?://www\\.gov\\.cn/xinwen/.*\\.html?")

	hfTitle := &service.HtmlFunc{
		QuerySelect: "#ti",
		F:           s.updateTitle,
	}

	hfContent := &service.HtmlFunc{
		QuerySelect: ".pages_content",
		F:           s.updateContent,
	}

	return service.NormalRule(rule, hfTitle, hfContent)
}

func (s *ScienceContentColly) chinataxCollector() *service.Rule {

	rule := regexp.MustCompile("https?://www\\.chinatax\\.gov\\.cn/.*\\.html?")

	hfTitle := &service.HtmlFunc{
		QuerySelect: "title",
		F:           s.updateTitle,
	}

	hfContent := &service.HtmlFunc{
		QuerySelect: "#fontzoom",
		F:           s.updateContent,
	}

	return service.NormalRule(rule, hfTitle, hfContent)
}
