package meta

import (
	"PolicySearchEngine/dao/database"
	"PolicySearchEngine/service"
	"PolicySearchEngine/utils"
	"errors"
	"fmt"
	"github.com/gocolly/colly"
)

const (
	initPage          = ""
	departmentID      = 1 // 科学技术部
	smallDepartmentID = 1 // 科学技术部
	provinceID        = 2 // 天津市
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
		colly.MaxDepth(1),
	)

	s.metaDal = &database.MetaDal{Db: database.MyDb()}
	s.dMapDal = &database.SmallDepartmentMapDal{Db: database.MyDb()}
}

func (s *ScienceMetaColly) PageTraverse() {
	// todo 根据initPage起始页，确定要遍历的页数，暂时写死，等待后续优化
	s.startPages = append(s.startPages, //initPage,
		//"https://kxjs.tj.gov.cn/ZWGK4143/TZGG2079/index.html",
		//"https://kxjs.tj.gov.cn/ZWGK4143/TZGG2079/index_1.html",
		//"https://kxjs.tj.gov.cn/ZWGK4143/TZGG2079/index_2.html",
		//"https://kxjs.tj.gov.cn/ZWGK4143/TZGG2079/index_3.html",
		//"https://kxjs.tj.gov.cn/ZWGK4143/ZXGZ7816/KJJR3268/GZDT5932/",
		//"https://kxjs.tj.gov.cn/ZWGK4143/ZXGZ7816/KJJR3268/GZDT5932/index_1.html",
		//"https://kxjs.tj.gov.cn/ZWGK4143/ZXGZ7816/KJJR3268/GZDT5932/index_2.html",
		//"https://kxjs.tj.gov.cn/ZWGK4143/ZXGZ7816/KJJR3268/GZDT5932/index_3.html",
		//"https://kxjs.tj.gov.cn/ZWGK4143/JZBGHYT7229/index.html",
		//"https://kxjs.tj.gov.cn/ZWGK4143/JZBGHYT7229/index_1.html",
		//"https://kxjs.tj.gov.cn/ZWGK4143/JZBGHYT7229/index_2.html",
		//"https://kxjs.tj.gov.cn/ZWGK4143/JZBGHYT7229/index_3.html",
		//"https://kxjs.tj.gov.cn/ZWGK4143/KJGH20201119/",
		//"https://kxjs.tj.gov.cn/ZWGK4143/KJGH20201119/index_1.html",
		//"https://kxjs.tj.gov.cn/ZWGK4143/KJGH20201119/index_2.html",
		//"https://kxjs.tj.gov.cn/ZWGK4143/XZZF201019/fzzfjsndbg/",
		//"https://kxjs.tj.gov.cn/ZWGK4143/XZZF201019/fzzfjsgz/",
		"https://kxjs.tj.gov.cn/ZWGK4143/XZZF201019/fzzfjsgz/index_1.html",
		"https://kxjs.tj.gov.cn/ZWGK4143/XZZF201019/fzzfjsgz/index_2.html",
	)

}

func (s *ScienceMetaColly) Operate() {

	//redis.SetRedisStorage(s.c, "meta-sci", s.startPages)

	s.c.OnHTML("[class='news_list news_list2'] li", func(e *colly.HTMLElement) {

		url := e.Request.AbsoluteURL(e.ChildAttr("a", "href"))
		date := e.ChildText("span")
		title := e.ChildText("a")
		fmt.Println(date, title, url)

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
