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
	initPage          = "https://www.most.gov.cn/satp/kjzc/zh/index.html"
	departmentID      = 1  // 科学技术部
	smallDepartmentID = 1  // 科学技术部
	provinceID        = 35 // 中央
)

type ScienceMetaColly struct {
	c *colly.Collector
	// 遍历起始页
	startPages []string
	metaDal    *database.MetaDal
	dMapDal    *database.SmallDepartmentMapDal
}

func (s *ScienceMetaColly) Init() {

	s.c = colly.NewCollector(
		colly.AllowedDomains(
			"www.most.gov.cn",
			"www.gov.cn",
			"szs.mof.gov.cn",
			"www.chinatax.gov.cn",
		),
		colly.URLFilters(

			regexp.MustCompile("https?://www\\.most\\.gov\\.cn/xxgk/.*\\.html?"),
			regexp.MustCompile("https?://www\\.most\\.gov\\.cn/tztg/.*\\.html?"),
			regexp.MustCompile("https?://www\\.most\\.gov\\.cn/satp/kjzc/zh/.*\\.html?"),
			regexp.MustCompile("https?://www\\.most\\.gov\\.cn/kjbgz/.*\\.html?"),

			regexp.MustCompile("https?://szs\\.mof\\.gov\\.cn/zhengwuxinxi/zhengcefabu/.*\\.html?"),
			regexp.MustCompile("https?://www\\.chinatax\\.gov\\.cn/.*\\.html?"),
			regexp.MustCompile("https?://www\\.gov\\.cn/zhengce/content/.*\\.html?"),
			regexp.MustCompile("https?://www\\.gov\\.cn/xinwen/.*\\.html?"),

			regexp.MustCompile("https?://www\\.chinatax\\.gov\\.cn/.*\\.html?"),
			regexp.MustCompile("https?://www\\.gov\\.cn/gongbao/content/.*\\.htm"),
		),
		colly.DisallowedURLFilters(
			// 去除 404 页面
			regexp.MustCompile("http://www\\.mof\\.gov\\.cn/404\\.htm"),
		),
		colly.MaxDepth(1),
	)

	s.metaDal = &database.MetaDal{Db: database.MyDb()}
	s.dMapDal = &database.SmallDepartmentMapDal{Db: database.MyDb()}
}

func (s *ScienceMetaColly) PageTraverse() {
	// todo 根据initPage起始页，确定要遍历的页数，暂时写死，等待后续优化
	s.startPages = append(s.startPages,
		initPage,
		"https://www.most.gov.cn/satp/kjzc/zh/index_1.html",
		"https://www.most.gov.cn/satp/kjzc/zh/index_2.html",
		"https://www.most.gov.cn/satp/kjzc/zh/index_3.html",
		"https://www.most.gov.cn/satp/kjzc/zh/index_4.html",
		"https://www.most.gov.cn/satp/kjzc/zh/index_5.html",
		"https://www.most.gov.cn/satp/kjzc/zh/index_6.html",
	)
}

func (s *ScienceMetaColly) Operate() {

	//redis.SetRedisStorage(s.c, "meta-sci", s.startPages)

	s.c.OnHTML(".list-main ul li", func(e *colly.HTMLElement) {

		url := e.Request.AbsoluteURL(e.ChildAttr("a", "href"))
		date := e.ChildText("span")
		title := e.ChildText("a")

		err := s.c.Visit(url) //尝试访问获取到的URL
		if errors.Is(err, colly.ErrAlreadyVisited) {
			return
		} //如果 URL 已经被访问过，则跳过。
		if err != nil {
			fmt.Println(err.Error() + fmt.Sprintf(" %q -> %s\n", e.Text, url))
			return //如果访问过程中发生其他错误，打印错误信息并跳过。
		}

		dateTime, err := utils.StringToTime(date) //尝试将字符串格式的日期转换为time.Time类型。
		if err != nil {
			fmt.Println(err.Error() + fmt.Sprintf("Time Falted %s %q -> %s\n", date, title, url))
			return
		}

		metaID := s.metaDal.InsertMeta(dateTime, title, url, departmentID, provinceID) //将日期、标题、URL 等信息插入到元数据表中，并返回插入的元数据 ID。
		s.dMapDal.InsertDID(metaID, smallDepartmentID)                                 //将元数据 ID 和小部门 ID 插入到部门映射表中。

		fmt.Printf("Link found: %s %q -> %s\n\n", date, title, url)
	})

}

func (s *ScienceMetaColly) Run() {
	for _, page := range s.startPages {
		err := s.c.Visit(page)
		if err != nil {
			fmt.Println(fmt.Sprintf("page:%s, err:%+v", page, err))
		}
	}
}

func (s *ScienceMetaColly) Destroy() {
	// 下次运行是在一天后了，指向nil，保证内存释放，让gc自动去回收
	s.c = nil
	s.metaDal = nil
	s.startPages = nil
}

func (s *ScienceMetaColly) ExecuteWorkflow() {
	s.Init()
	s.PageTraverse()
	s.Operate()
	s.Run()
	s.Destroy()
}

var _ service.MetaCrawler = (*ScienceMetaColly)(nil)
