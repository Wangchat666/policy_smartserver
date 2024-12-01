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
	initPage          = "https://kjt.ln.gov.cn/kjt/kjzc/index.shtml"
	departmentID      = 1 // 科学技术部
	smallDepartmentID = 1 // 科学技术部
	provinceID        = 6 // 辽宁省
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
	//s.startPages = append(s.startPages, //initPage,
	//"https://kjt.ln.gov.cn/kjt/kjzc/lnkjzc/index.shtml",
	//"https://kjt.ln.gov.cn/kjt/kjzc/lnkjzc/5b491c27-2.shtml",
	//"https://kjt.ln.gov.cn/kjt/kjzc/lnkjzc/5b491c27-3.shtml",
	//"https://kjt.ln.gov.cn/kjt/kjzc/lnkjzc/5b491c27-4.shtml",
	//"https://kjt.ln.gov.cn/kjt/kjzc/dfkjzc/index.shtml",
	//"https://kjt.ln.gov.cn/kjt/kjzc/dfkjzc/9003f6a4-2.shtml",
	//"https://kjt.ln.gov.cn/kjt/kjzc/gjkjzc/index.shtml",
	//"https://kjt.ln.gov.cn/kjt/kjzc/gjkjzc/31b346a5-2.shtml",
	//"https://kjt.ln.gov.cn/kjt/kjzc/zcjd/index.shtml"，

	//农业科技
	//for i := 1; i < 9; i++ {
	//	s.startPages = append(s.startPages, fmt.Sprintf("https://kjt.ln.gov.cn/kjt/kjgz/nykj/07a2b7b3-%d.shtml", i))
	//}
	//社发科技
	//for i := 1; i < 11; i++ {
	//	s.startPages = append(s.startPages, fmt.Sprintf("https://kjt.ln.gov.cn/kjt/kjgz/sfkj/8a840d54-%d.shtml", i))
	//}
	//科技人才国际合作
	//for i := 1; i < 8; i++ {
	//	s.startPages = append(s.startPages, fmt.Sprintf("https://kjt.ln.gov.cn/kjt/kjgz/kjrc/d144422b-%d.shtml", i))
	//}
	//成果转化
	//for i := 1; i < 9; i++ {
	//	s.startPages = append(s.startPages, fmt.Sprintf("https://kjt.ln.gov.cn/kjt/kjgz/cgzh/9592dc11-%d.shtml", i))
	//}
	//科研诚信
	//for i := 1; i < 11; i++ {
	//	s.startPages = append(s.startPages, fmt.Sprintf("https://kjt.ln.gov.cn/kjt/kjgz/kjhz/a401de05-%d.shtml", i))
	//}
	//政策体改
	//for i := 1; i < 11; i++ {
	//	s.startPages = append(s.startPages, fmt.Sprintf("https://kjt.ln.gov.cn/kjt/kjgz/tzgg/9cc94703-%d.shtml", i))
	//}
	//科技企业
	//for i := 1; i < 10; i++ {
	//	s.startPages = append(s.startPages, fmt.Sprintf("https://kjt.ln.gov.cn/kjt/kjgz/kjfw/1b34d93e-%d.shtml", i))
	//}
	//高新科技
	//for i := 1; i < 11; i++ {
	//	s.startPages = append(s.startPages, fmt.Sprintf("https://kjt.ln.gov.cn/kjt/kjgz/gxkj/ca44553e-%d.shtml", i))
	//}
	//科技计划
	//for i := 1; i < 4; i++ {
	//	s.startPages = append(s.startPages, fmt.Sprintf("https://kjt.ln.gov.cn/kjt/kjgz/kjjh/5cd14ef9-%d.shtml", i))
	//}
	//创新高地
	//for i := 1; i < 8; i++ {
	//	s.startPages = append(s.startPages, fmt.Sprintf("https://kjt.ln.gov.cn/kjt/kjgz/cxjd/94a0712b-%d.shtml", i))
	//}
	//党建人事
	//for i := 1; i < 9; i++ {
	//	s.startPages = append(s.startPages, fmt.Sprintf("https://kjt.ln.gov.cn/kjt/kjgz/djgz/621ce51b-%d.shtml", i))
	//}
	//信息技术
	//for i := 1; i < 8; i++ {
	//	s.startPages = append(s.startPages, fmt.Sprintf("https://kjt.ln.gov.cn/kjt/kjgz/djgz_136326/60e1fd81-%d.shtml", i))
	//}
	//医药科技
	//for i := 1; i < 5; i++ {
	//	s.startPages = append(s.startPages, fmt.Sprintf("https://kjt.ln.gov.cn/kjt/kjgz/djgz_136327/390058b6-%d.shtml", i))
	//}
	//资金管理
	s.startPages = append(s.startPages, fmt.Sprintf("https://kjt.ln.gov.cn/kjt/kjgz/zjgl/index.shtml"))
	//辽宁实验室
	//for i := 1; i < 3; i++ {
	//	s.startPages = append(s.startPages, fmt.Sprintf("https://kjt.ln.gov.cn/kjt/kjgz/lnsys/8a840d54-%d.shtml", i))
	//}
}

func (s *ScienceMetaColly) Operate() {

	//redis.SetRedisStorage(s.c, "meta-sci", s.startPages)

	s.c.OnHTML(".govCListBox li", func(e *colly.HTMLElement) {

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
