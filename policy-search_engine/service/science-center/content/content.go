package content

import (
	"PolicySearchEngine/dao/database"
	"PolicySearchEngine/model"
	"PolicySearchEngine/service"
	"fmt"
)

const (
	departmentID = 1  // 科学技术部
	provinceID   = 35 // 中央
)

type ScienceContentColly struct {
	rules []*service.Rule
	// 等待处理的url队列
	waitQueue  *[]model.Meta
	metaDal    *database.MetaDal
	contentDal *database.ContentDal
	dMapDal    *database.SmallDepartmentMapDal
}

func (s *ScienceContentColly) Init() {
	// 注册规则
	s.rules = append(s.rules, s.getRules()...)         //获取爬取规则
	s.metaDal = &database.MetaDal{Db: database.MyDb()} //初始化 s.metaDal 成员变量为一个新的 database.MetaDal 实例，并将其 Db 字段设置为 database.MyDb() 返回的数据库连接实例。
	s.contentDal = &database.ContentDal{Db: database.MyDb()}
	s.dMapDal = &database.SmallDepartmentMapDal{Db: database.MyDb()}
}

// Import 分批次导入
func (s *ScienceContentColly) Import() (success bool) {
	// todo 1. 暂时全量导入 2. 监控时需要单独区分哪些读过
	metaList := s.metaDal.GetAllMeta(departmentID, provinceID) //传入 departmentID 和 provinceID 参数，获取所有元数据的列表。
	if metaList == nil || len(*metaList) == 0 {
		return false
	}
	s.waitQueue = metaList
	return true
} //从数据库中获取所有元数据，并将其存储到 ScienceContentColly 实例的等待队列中。

func (s *ScienceContentColly) Run() {

	dealMeta := func(meta *model.Meta) {
		var match bool
		for _, rule := range s.rules {
			if rule.R.MatchString(meta.Url) { //检查当前规则是否匹配元数据的 URL。
				if err := rule.C.Visit(meta.Url); err != nil {
					fmt.Println(err)
				}
				match = true
				break
			}
		}
		if !match {
			fmt.Printf("url:%s 未匹配到任何规则\n", meta.Url)
		}
	}

	for _, meta := range *s.waitQueue { //遍历 s.waitQueue 中的所有元数据
		fmt.Printf("I'm dealing %s...\n", meta.Url)
		dealMeta(&meta)
	}
}

func (s *ScienceContentColly) Destroy() {
	s.rules = nil
	s.metaDal = nil
	s.contentDal = nil
	s.waitQueue = nil
}

func (s *ScienceContentColly) ExecuteWorkflow() {
	s.Init()
	if s.Import() {
		s.Run()
	}
	s.Destroy()
}

var _ service.ContentCrawler = (*ScienceContentColly)(nil)
