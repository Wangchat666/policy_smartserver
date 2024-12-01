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
	"strings"
)

const (
	departmentID      = 91 // 国务院
	smallDepartmentID = 4  // 国务院文件
	provinceID        = 35 // 中央
)

type Resp struct {
	Code     int    `json:"code"`
	Msg      string `json:"msg"`
	SearchVO struct {
		TotalCount  int `json:"totalCount"`
		CurrentPage int `json:"currentPage"`
		Totalpage   int `json:"totalpage"`
		ListVO      []struct {
			Title      string `json:"title"`
			PubtimeStr string `json:"pubtimeStr"`
			URL        string `json:"url"`
			Puborg     string `json:"puborg"`
		} `json:"listVO"`
	} `json:"searchVO"`
}

var smallDepartmentIDMap = map[string]uint{
	"科学技术部":      1,
	"教育部":        2,
	"工业和信息化部":    3,
	"国务院文件":      4,
	"外交部":        5,
	"国家发展和改革委员会": 6,
	"国家民族事务委员会":  7,
	"公安部":        8,
	"国家安全部":      9,
	"民政部":        10,
	"司法部":        11,
	"财政部":        12,
	"人力资源和社会保障部": 13,
	"自然资源部":      14,
	"生态环境部":      15,
	"住房和城乡建设部":   16,
	"交通运输部":      17,
	"水利部":        18,
	"农业农村部":      19,
	"商务部":        20,
	"文化和旅游部":     21,
	"国家卫生健康委员会":  22,
	"退役军人事务部":    23,
	"应急管理部":      24,
	"人民银行":       25,
	"审计署":        26,
	"国务院国有资产监督管理委员会": 27,
	"海关总署":          28,
	"国家税务总局":        29,
	"国家市场监督管理总局":    30,
	"国家金融监督管理总局":    31,
	"国家广播电视总局":      32,
	"国家体育总局":        33,
	"国家统计局":         34,
	"国家国际发展合作署":     35,
	"国家医疗保障局":       36,
	"国家机关事务管理局":     37,
	"国家标准化管理委员会":    38,
	"国家新闻出版署":       39,
	"国家版权局":         40,
	"国家互联网信息办公室":    41,
	"中国科学院":         42,
	"中国社会科学院":       43,
	"中国工程院":         44,
	"中国气象局":         45,
	"中国银行保险监督管理委员会": 46,
	"中国证券监督管理委员会":   47,
	"国家信访局":         48,
	"国家粮食和物资储备局":    49,
	"国家能源局":         50,
	"国家数据局":         51,
	"国家国防科技工业局":     52,
	"国家烟草专卖局":       53,
	"国家移民管理局":       54,
	"国家林业和草原局":      55,
	"国家铁路局":         56,
	"中国民用航空局":       57,
	"国家邮政局":         58,
	"国家文物局":         59,
	"国家中医药管理局":      60,
	"国家疾病预防控制局":     61,
	"国家矿山安全监察局":     62,
	"国家消防救援局":       63,
	"国家外汇管理局":       64,
	"国家药品监督管理局":     65,
	"国家知识产权局":       66,
	"国家公务员局":        67,
	"国家档案局":         68,
	"国家保密局":         69,
	"国家密码管理局":       70,
	"国家航天局":         71,
	"国家原子能机构":       72,
	"国家宗教事务局":       73,
	"国务院台湾事务办公室":    74,
	"国家乡村振兴局":       75,
	"国家核安全局":        76,
	"国家认证认可监督管理委员会": 77,
	"国家语言文字工作委员会":   78,
	"国家电影局":         79,
}

type StateCouncilMetaColly struct {
	c *colly.Collector
	// 遍历起始页
	startPages  []string
	currentPage string
	metaDal     *database.MetaDal
	dMapDal     *database.SmallDepartmentMapDal
}

func (s *StateCouncilMetaColly) Init() {
	s.metaDal = &database.MetaDal{Db: database.MyDb()}
	s.dMapDal = &database.SmallDepartmentMapDal{Db: database.MyDb()}
}

func (s *StateCouncilMetaColly) PageTraverse() {
	// 国务院下设部门
	bmInitPage := "https://sousuo.www.gov.cn/search-gov/data?t=zhengcelibrary_bm&sort=score&sortType=1&searchfield=title&n=100&type=gwyzcwjk&p="
	//for i := 1; i < 112; i++ {
	for i := 74; i < 112; i++ {
		s.startPages = append(s.startPages, bmInitPage+strconv.Itoa(i))
	}

	// 国务院文件
	gwInitPage := "https://sousuo.www.gov.cn/search-gov/data?" +
		"t=zhengcelibrary_gw&sort=score&sortType=1&searchfield=title&n=100&type=gwyzcwjk&p="
	for i := 1; i < 61; i++ {
		s.startPages = append(s.startPages, gwInitPage+strconv.Itoa(i))
	}

}

func (s *StateCouncilMetaColly) Operate() {

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
		fmt.Printf("读取响应体时发生错误：%s\n", err)
		return
	}

	// 解析JSON
	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		fmt.Printf("解析JSON时发生错误：%s\n", err)
		fmt.Println("当前爬取页面为：" + s.currentPage)
		return
	}

	var resp Resp

	// 解析JSON
	if err := json.Unmarshal(body, &resp); err != nil {
		fmt.Printf("解析JSON时发生错误：%s\n", err)
		fmt.Println("当前爬取页面为：" + s.currentPage)
		return
	}

	for _, list := range resp.SearchVO.ListVO {

		dateTime, _ := utils.StringToTimeByDot(list.PubtimeStr)
		metaID := s.metaDal.InsertMeta(dateTime, list.Title, list.URL, departmentID, provinceID)

		// 如果url中有zhengcelibrary_bm代表下设部门，根据puborg查询对应部门，否则直接国务院
		if strings.Contains(s.currentPage, "zhengcelibrary_bm") {
			pubs := strings.Fields(list.Puborg)
			for _, pub := range pubs {
				sdID := smallDepartmentIDMap[pub]

				if sdID == 0 {
					fmt.Println(list.Puborg)
				}

				s.dMapDal.InsertDID(metaID, sdID)
			}
		} else {
			s.dMapDal.InsertDID(metaID, smallDepartmentID)
		}

		fmt.Printf("标题：%s，URL：%s，日期：%s\n", list.Title, list.URL, dateTime)
	}

	//<-time.After(3 * time.Second)

}

func (s *StateCouncilMetaColly) Run() {
	for _, page := range s.startPages {
		s.currentPage = page
		s.Operate()
	}
}

func (s *StateCouncilMetaColly) Destroy() {
	// 下次运行是在一天后了，指向nil，保证内存释放，让gc自动去回收
	s.c = nil
	s.metaDal = nil
	s.startPages = nil
}

func (s *StateCouncilMetaColly) ExecuteWorkflow() {
	s.Init()
	s.PageTraverse()
	s.Run()
	s.Destroy()
}

var _ service.MetaCrawler = (*StateCouncilMetaColly)(nil)
