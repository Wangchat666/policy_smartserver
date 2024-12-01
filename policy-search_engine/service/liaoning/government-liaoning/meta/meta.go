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
	initPage     = ""
	departmentID = 92 // 人民政府
	provinceID   = 6  // 辽宁省
)

var smallDepartmentIDMap = map[string]uint{
	"辽宁省人民政府": 81,
	"沈阳市人民政府": 82,
	"大连市人民政府": 83,
	"大连市统计局":  84,
}

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
	//辽政发
	//for i := 1; i <= 42; i++ {
	//	s.startPages = append(s.startPages, fmt.Sprintf("https://www.ln.gov.cn/web/zwgkx/zfwj/605dc429-%d.shtml", i))
	//}
	//省政府令
	//for i := 1; i <= 9; i++ {
	//	s.startPages = append(s.startPages, fmt.Sprintf("https://www.ln.gov.cn/web/zwgkx/zfwj/szfl/661e4dc5-%d.shtml", i))
	//}
	//辽政
	//for i := 1; i <= 2; i++ {
	//	s.startPages = append(s.startPages, fmt.Sprintf("https://www.ln.gov.cn/web/zwgkx/zfwj/lz/b0e30107-%d.shtml", i))
	//}
	//辽政办发 标签有变化  ul[class='list-ul-one list-ul-two' 其余标签相同
	//for i := 1; i <= 68; i++ {
	//	s.startPages = append(s.startPages, fmt.Sprintf("https://www.ln.gov.cn/web/zwgkx/zfwj/szfbgtwj/3c92e3a6-%d.shtml", i))
	//}
	//辽政办
	//for i := 1; i <= 2; i++ {
	//	s.startPages = append(s.startPages, fmt.Sprintf("https://www.ln.gov.cn/web/zwgkx/zfwj/lzb/56a176e0-%d.shtml", i))
	//}
	//各市文件  沈阳市
	//s.startPages = append(s.startPages, "https://www.shenyang.gov.cn/zwgk/zcwj/fggz/sdfxfgzl/index.html",
	//	"https://www.shenyang.gov.cn/zwgk/zcwj/fggz/sdfxfgzl/index_1.html",
	//	"https://www.shenyang.gov.cn/zwgk/zcwj/fggz/sdfxfgzl/index_2.html",
	//)
	//s.startPages = append(s.startPages, "https://www.shenyang.gov.cn/zwgk/zcwj/xzgfxwj/szfgfxwj1/")
	//大连市
	//for i := 1; i <= 52; i++ {
	//	s.startPages = append(s.startPages, fmt.Sprintf("https://www.dl.gov.cn/module/xxgk/search.jsp?divid=div4&jdid=1&area=&infotypeId=DL00302&isAllList=1&standardXxgk=2&currpage=%d", i))
	//}
	for i := 1; i <= 1; i++ {
		s.startPages = append(s.startPages, fmt.Sprintf("https://www.dl.gov.cn/module/xxgk/search.jsp?divid=div4&jdid=1&area=&infotypeId=DL0030203&isAllList=1&standardXxgk=2&currpage=%d", i))
	}
	//for i := 1; i <= 3; i++ {
	//	s.startPages = append(s.startPages, fmt.Sprintf("https://www.dl.gov.cn/module/xxgk/search.jsp?divid=div4&jdid=1&area=&infotypeId=DL01205&isAllList=1&standardXxgk=2&currpage=%d", i))
	//}

}

func (s *ScienceMetaColly) Operate() {

	//redis.SetRedisStorage(s.c, "meta-sci", s.startPages)

	s.c.OnHTML("ul[class='list-ul-one list-ul-two hy']", func(e *colly.HTMLElement) {

		url := e.Request.AbsoluteURL(e.ChildAttr("a#TITLETEXT", "href"))
		title := e.ChildText("a#TITLETEXT")
		date := e.ChildText("a.subStr2")[:10]
		//fmt.Println(url, title, date)

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
		s.dMapDal.InsertDID(metaID, smallDepartmentIDMap["辽宁省人民政府"])

		fmt.Printf("Link found: %s %q -> %s\n\n", date, title, url)
	})
	s.c.OnHTML("ul[class='xxgk_rul'] li", func(e *colly.HTMLElement) {

		url := e.ChildAttr("a", "onclick")
		url = regexp.MustCompile(`https?://\S+[.][html|si]+`).FindString(url)
		title := e.ChildText("a")
		date := e.ChildText("span")
		//fmt.Println(url, title, date)

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
		s.dMapDal.InsertDID(metaID, smallDepartmentIDMap["沈阳市人民政府"])

		fmt.Printf("Link found: %s %q -> %s\n\n", date, title, url)
	})
	s.c.OnHTML("ul[class='list-sp'] ", func(e *colly.HTMLElement) {

		url := e.ChildAttr("[class='title']", "href")
		title := e.ChildText("[class='title']")
		date := e.ChildText("[class='time_pub']")
		//fmt.Println(url, title, date)

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
		s.dMapDal.InsertDID(metaID, smallDepartmentIDMap["沈阳市人民政府"])

		fmt.Printf("Link found: %s %q -> %s\n\n", date, title, url)
	})
	s.c.OnHTML("[style='font-size: 9pt;color: #3D3D3D;font-family: 微软雅黑;line-height: 180%;']", func(e *colly.HTMLElement) {

		url := e.Request.AbsoluteURL(e.ChildAttr("a", "href"))
		title := e.ChildAttr("a", "title")
		date := e.ChildText("[width='80']")
		//fmt.Println(url, title, date)

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
		s.dMapDal.InsertDID(metaID, smallDepartmentIDMap["大连市人民政府"])

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
