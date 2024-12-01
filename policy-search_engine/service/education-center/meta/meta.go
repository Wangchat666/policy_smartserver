package meta

import (
	"PolicySearchEngine/dao/database"
	"PolicySearchEngine/service"
	"PolicySearchEngine/utils"
	"errors"
	"fmt"
	"github.com/gocolly/colly"
	"regexp"
)

const (
	initPage          = "http://www.moe.gov.cn/jyb_sjzl/sjzl_zcfg/moe_418/index.html"
	departmentID      = 2  // 教育部
	smallDepartmentID = 2  // 教育部
	provinceID        = 35 // 中央
)

type EducationMetaColly struct {
	c *colly.Collector
	// 遍历起始页
	startPages []string
	metaDal    *database.MetaDal
	dMapDal    *database.SmallDepartmentMapDal
}

func (s *EducationMetaColly) Init() {

	s.c = colly.NewCollector(
		colly.AllowedDomains(
			"www.moe.gov.cn",
		),
		colly.URLFilters(
			regexp.MustCompile("https?://www\\.moe\\.gov\\.cn/jyb_sjzl/sjzl_zcfg/moe_418/.*\\.html?"),
			regexp.MustCompile("https?://www\\.moe\\.gov\\.cn/jyb_sjzl/sjzl_zcfg/zcfg_jyfl/.*\\.html?"),
			regexp.MustCompile("https?://www\\.moe\\.gov\\.cn/jyb_sjzl/sjzl_zcfg/zcfg_jyxzfg/.*\\.html?"),
			regexp.MustCompile("https?://www\\.moe\\.gov\\.cn/jyb_sjzl/sjzl_zcfg/zcfg_qtxgfl/.*\\.html?"),
			regexp.MustCompile("https?://www\\.moe\\.gov\\.cn/jyb_sjzl/sjzl_zcfg/zcfg_qtxgfg/.*\\.html?"),
			regexp.MustCompile("https?://www\\.moe\\.gov\\.cn/jyb_sjzl/sjzl_zcfg/zcfg_jybmgz/.*\\.html?"),
			regexp.MustCompile("https?://www\\.moe\\.gov\\.cn/jyb_sjzl/sjzl_zcfg/zcfg_gdxxzch/.*\\.html?"),
			regexp.MustCompile("https?://www\\.moe\\.gov\\.cn/srcsite/.*\\.html?"),
		),
		colly.MaxDepth(1),
	)

	s.metaDal = &database.MetaDal{Db: database.MyDb()}
	s.dMapDal = &database.SmallDepartmentMapDal{Db: database.MyDb()}
}

func (s *EducationMetaColly) PageTraverse() {

	page := []string{
		"http://www.moe.gov.cn/jyb_sjzl/sjzl_zcfg/moe_418/index.html",     // 宪法
		"http://www.moe.gov.cn/jyb_sjzl/sjzl_zcfg/zcfg_jyfl/index.html",   // 教育法规
		"http://www.moe.gov.cn/jyb_sjzl/sjzl_zcfg/zcfg_jyxzfg/index.html", // 教育行政法规
		"http://www.moe.gov.cn/jyb_sjzl/sjzl_zcfg/zcfg_qtxgfl/index.html", //其他相关法律
		"http://www.moe.gov.cn/jyb_sjzl/sjzl_zcfg/zcfg_qtxgfg/index.html", //其他相关法规
	}
	page1 := []string{ // 教育部门规章
		"http://www.moe.gov.cn/jyb_sjzl/sjzl_zcfg/zcfg_jybmgz/index.html",
		"http://www.moe.gov.cn/jyb_sjzl/sjzl_zcfg/zcfg_jybmgz/index_1.html",
		"http://www.moe.gov.cn/jyb_sjzl/sjzl_zcfg/zcfg_jybmgz/index_2.html",
	}
	page2 := []string{ // 高等学校章程
		"http://www.moe.gov.cn/jyb_sjzl/sjzl_zcfg/zcfg_gdxxzch/index.html",
		"http://www.moe.gov.cn/jyb_sjzl/sjzl_zcfg/zcfg_gdxxzch/index_1.html",
		"http://www.moe.gov.cn/jyb_sjzl/sjzl_zcfg/zcfg_gdxxzch/index_2.html",
		"http://www.moe.gov.cn/jyb_sjzl/sjzl_zcfg/zcfg_gdxxzch/index_3.html",
		"http://www.moe.gov.cn/jyb_sjzl/sjzl_zcfg/zcfg_gdxxzch/index_4.html",
		"http://www.moe.gov.cn/jyb_sjzl/sjzl_zcfg/zcfg_gdxxzch/index_5.html",
		"http://www.moe.gov.cn/jyb_sjzl/sjzl_zcfg/zcfg_gdxxzch/index_6.html",
		"http://www.moe.gov.cn/jyb_sjzl/sjzl_zcfg/zcfg_gdxxzch/index_7.html",
		"http://www.moe.gov.cn/jyb_sjzl/sjzl_zcfg/zcfg_gdxxzch/index_8.html",
		"http://www.moe.gov.cn/jyb_sjzl/sjzl_zcfg/zcfg_gdxxzch/index_9.html",
	}

	s.startPages = append(s.startPages, initPage)
	s.startPages = append(s.startPages, page...)
	s.startPages = append(s.startPages, page1...)
	s.startPages = append(s.startPages, page2...)
}

func (s *EducationMetaColly) Operate() {

	//redis.SetRedisStorage(s.c, "meta-edu", s.startPages)

	s.c.OnHTML(".moe-list ul li", func(e *colly.HTMLElement) {

		url := e.Request.AbsoluteURL(e.ChildAttr("a", "href"))
		date := e.ChildText("span")
		title := e.ChildText("a")

		err := s.c.Visit(url)
		if errors.Is(err, colly.ErrAlreadyVisited) {
			return
		}
		if err != nil {
			fmt.Println(err.Error() + fmt.Sprintf(" %q -> %s\n", e.Text, url))
			return
		}

		dateTime, err := utils.StringToTime(date)
		if err != nil {
			fmt.Println(err.Error() + fmt.Sprintf("Time Falted %s %q -> %s\n", date, title, url))
			return
		}

		metaID := s.metaDal.InsertMeta(dateTime, title, url, departmentID, provinceID)
		s.dMapDal.InsertDID(metaID, smallDepartmentID)

		fmt.Printf("Link found: %s %q -> %s\n\n", date, title, url)
	})

}

func (s *EducationMetaColly) Run() {
	for _, page := range s.startPages {
		err := s.c.Visit(page)
		if err != nil {
			fmt.Println(fmt.Sprintf("page:%s, err:%+v", page, err))
		}
	}
}

func (s *EducationMetaColly) Destroy() {
	// 下次运行是在一天后了，指向nil，保证内存释放，让gc自动去回收
	s.c = nil
	s.metaDal = nil
	s.startPages = nil
}

func (s *EducationMetaColly) ExecuteWorkflow() {
	s.Init()
	s.PageTraverse()
	s.Operate()
	s.Run()
	s.Destroy()
}

var _ service.MetaCrawler = (*EducationMetaColly)(nil)
