package content

import (
	"PolicySearchEngine/dao/database"
	"PolicySearchEngine/model"
	"PolicySearchEngine/service"
	"fmt"
)

const (
	departmentID = 2  // 教育部
	provinceID   = 35 // 中央
)

type EducationContentColly struct {
	rules []*service.Rule
	// 等待处理的url队列
	waitQueue  *[]model.Meta
	metaDal    *database.MetaDal
	contentDal *database.ContentDal
	dMapDal    *database.SmallDepartmentMapDal
}

func (s *EducationContentColly) Init() {
	// 注册规则
	s.rules = append(s.rules, s.getRules()...)
	s.metaDal = &database.MetaDal{Db: database.MyDb()}
	s.contentDal = &database.ContentDal{Db: database.MyDb()}
	s.dMapDal = &database.SmallDepartmentMapDal{Db: database.MyDb()}
}

// Import 分批次导入
func (s *EducationContentColly) Import() (success bool) {
	// todo 1. 暂时全量导入 2. 监控时需要单独区分哪些读过
	metaList := s.metaDal.GetAllMeta(departmentID, provinceID)
	if metaList == nil || len(*metaList) == 0 {
		return false
	}
	s.waitQueue = metaList
	return true
}

func (s *EducationContentColly) Run() {

	dealMeta := func(meta *model.Meta) {
		var match bool
		for _, rule := range s.rules {
			if rule.R.MatchString(meta.Url) {
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

	for _, meta := range *s.waitQueue {
		fmt.Printf("I'm dealing %s...\n", meta.Url)
		dealMeta(&meta)
	}
}

func (s *EducationContentColly) Destroy() {
	s.rules = nil
	s.metaDal = nil
	s.contentDal = nil
	s.waitQueue = nil
}

func (s *EducationContentColly) ExecuteWorkflow() {
	s.Init()
	if s.Import() {
		s.Run()
	}
	s.Destroy()
}

var _ service.ContentCrawler = (*EducationContentColly)(nil)
