package meta

import (
	"PolicySearchEngine/dao/database"
	"PolicySearchEngine/service"
	"PolicySearchEngine/utils"
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly"
	"io/ioutil"
	"net/http"
	"strconv"
)

const (
	departmentID      = 3  // 工信部
	smallDepartmentID = 3  // 工信部
	provinceID        = 35 // 中央
)

type IndustryInformatizationMetaColly struct {
	c *colly.Collector
	// 遍历起始页
	startPages  []string
	currentPage string
	metaDal     *database.MetaDal
	dMapDal     *database.SmallDepartmentMapDal
}

func (s *IndustryInformatizationMetaColly) Init() {
	s.metaDal = &database.MetaDal{Db: database.MyDb()}
	s.dMapDal = &database.SmallDepartmentMapDal{Db: database.MyDb()}
}

func (s *IndustryInformatizationMetaColly) PageTraverse() {
	initPage1 := "https://wap.miit.gov.cn/search-front-server/api/search/info?cateid=59&pg=10&selectFields=title,url,cdate&p="
	for i := 1; i < 18; i++ {
		s.startPages = append(s.startPages, initPage1+strconv.Itoa(i))
	}

	initPage2 := "https://wap.miit.gov.cn/search-front-server/api/search/info?cateid=61&pg=10&selectFields=title,url,cdate&p="
	for i := 1; i < 73; i++ {
		s.startPages = append(s.startPages, initPage2+strconv.Itoa(i))
	}
}

func (s *IndustryInformatizationMetaColly) Operate() {

	// 发送GET请求
	response, err := http.Get(s.currentPage)
	if err != nil {
		fmt.Printf("发生错误：%s\n", err)
		return
	}
	defer response.Body.Close()

	// 读取响应体
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("读取响应体时发生错误：%s\n, body: %v", err, body)
		return
	}

	// 解析JSON
	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		fmt.Printf("解析JSON时发生错误：%s\n, data: %v", err, data)
		return
	}

	var resp struct {
		Data struct {
			SearchResult struct {
				DataResults []struct {
					Data struct {
						Title       string `json:"title"`
						URL         string `json:"url"`
						JsearchDate string `json:"jsearch_date"`
					} `json:"data"`
				} `json:"dataResults"`
			} `json:"searchResult"`
		} `json:"data"`
	}

	// 解析JSON
	if err := json.Unmarshal(body, &resp); err != nil {
		fmt.Printf("解析JSON时发生错误：%s\n", err)
		return
	}

	// 提取所需字段
	for _, result := range resp.Data.SearchResult.DataResults {
		dateTime, _ := utils.StringToTime(result.Data.JsearchDate)
		url := "https://wap.miit.gov.cn" + result.Data.URL

		metaID := s.metaDal.InsertMeta(dateTime, result.Data.Title, url, departmentID, provinceID)
		s.dMapDal.InsertDID(metaID, smallDepartmentID)

		fmt.Printf("标题：%s，URL：%s，日期：%s\n", result.Data.Title, url, result.Data.JsearchDate)
	}

}

func (s *IndustryInformatizationMetaColly) Run() {
	for _, page := range s.startPages {
		s.currentPage = page
		s.Operate()
	}
}

func (s *IndustryInformatizationMetaColly) Destroy() {
	// 下次运行是在一天后了，指向nil，保证内存释放，让gc自动去回收
	s.c = nil
	s.metaDal = nil
	s.startPages = nil
}

func (s *IndustryInformatizationMetaColly) ExecuteWorkflow() {
	s.Init()
	s.PageTraverse()
	s.Run()
	s.Destroy()
}

var _ service.MetaCrawler = (*IndustryInformatizationMetaColly)(nil)
