package meta

import (
	"PolicySearchEngine/dao/database"
	"PolicySearchEngine/service"
	"PolicySearchEngine/utils"
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly"
	"io"
	"net/http"
	"strconv"
)

const (
	provinceID = 35 // 中央
)

type Resp struct {
	Result struct {
		DepartmentName         string `json:"departmentName"`
		DepartmentID           string `json:"departmentId"`
		DepartmentPolicyNotice struct {
			Records []struct {
				Title       string `json:"title"`
				ReleaseTime string `json:"releaseTime"`
				EntranceURL string `json:"entranceUrl"`
			} `json:"records"`
		} `json:"departmentPolicyNotice"`
	} `json:"result"`
}

var departmentRequestIDMap = map[string]string{
	"国家药监局":         "1351878865153187842",
	"国家商务部":         "1351879097299525634",
	"国家艺术基金":        "1351879343924600834",
	"国家发改委":         "1351880206093148161",
	"国家卫生健康委员会":     "1351881180262195202",
	"国家医疗保障局":       "1351881229364912129",
	"国家知识产权局":       "1351882518895288322",
	"国家科信平台":        "1351882541330620417",
	"高新技术企业认定工作网":   "1351882755709886465",
	"国家农业农村部":       "1405353069525409794",
	"国家生态环境部":       "1405353069928062977",
	"国家退役军人事务部":     "1405353070360076289",
	"国家交通运输部":       "1405353070519459841",
	"国家文化和旅游部":      "1405353070880169986",
	"国家水利部":         "1405353071266045953",
	"国家人力资源和社会保障部":  "1405353073027653633",
	"国家林业和草原局":      "1405360634548264961",
	"国家行政学院":        "1405360634896392194",
	"国家信访局":         "1405360634963501057",
	"国家档案局":         "1405360635122884610",
	"国家统计局":         "1405360635441651714",
	"国家粮食和物资储备局":    "1405360635638784002",
	"国家版权局":         "1405360636146294786",
	"国家中医药管理局":      "1405360636339232769",
	"国家能源局":         "1405360636905463810",
	"国家移民管理局":       "1405360637605912578",
	"国家外汇管理局":       "1405360638100840449",
	"中国工程院":         "1405360638574796801",
	"国家审计署":         "1405360638725791746",
	"国务院港澳事务办公室":    "1405360639300411393",
	"国家应急管理部":       "1405360639472377858",
	"国务院参事室":        "1405360639916974082",
	"国务院台湾事务办公室":    "1405360640374153217",
	"国家广播电视总局":      "1405360641594695681",
	"国家原子能机构":       "1405360641884102658",
	"国家邮政局":         "1405360642597134337",
	"国家机关事务管理局":     "1405360643385663490",
	"中国气象局":         "1405360644778172417",
	"国家认证认可监督管理委员会": "1405360645277294594",
	"国家体育总局":        "1405360645470232578",
	"国家密码管理局":       "1405360646007103489",
	"国家标准化管理委员会":    "1405360646208430082",
	"中国银行保险监督管理委员会": "1405715523652550658",
	"国家自然资源部":       "1422589309366685697",
	"国家财政部":         "1422589309450571778",
	"国家能源局南方监管局":    "1550406605681520641",
	"国家住房和城乡建设部":    "1550414117268942849",
	"国家药品监督管理局食品药品审核查验中心": "1684112595684974594",
	"中国食品药品检定研究院":         "1684112596922294273",
	"国家组织医用耗材联合采购平台":      "1684112598264471553",
	"国家市场监督管理总局":          "1684112599279493122",
	"国家民政部":               "1684112600558755841",
	"国家药典委员会":             "1684459886858018818",
	"国家税务总局":              "1685909443408035841",
	"国家药品监督管理局医疗器械技术审评中心": "1689216239644225537",
	"国家海关总署":              "1689216244824190977",
}

// 要爬取的部门
var departmentIDMap = map[string]uint{
	"国家发改委":        14,
	"国家民政部":        21,
	"国家财政部":        23,
	"国家人力资源和社会保障部": 24,
	"国家生态环境部":      26,
	"国家交通运输部":      28,
	"国家水利部":        29,
	"国家农业农村部":      30,
	"国家商务部":        31,
	"国家文化和旅游部":     32,
	"国家卫生健康委员会":    33,
	"国家退役军人事务部":    34,
	"国家海关总署":       39,
	"国家税务总局":       40,
	"国家市场监督管理总局":   41,
	"国家统计局":        45,
	"国家医疗保障局":      47,
	"国家版权局":        51,
	"中国工程院":        55,
	"国家信访局":        59,
	"国家粮食和物资储备局":   60,
	"国家能源局":        61,
	"国家移民管理局":      65,
	"国家林业和草原局":     66,
	"国家邮政局":        69,
	"国家中医药管理局":     71,
	"国家外汇管理":       75,
	"国家药监局":        76,
	"国家知识产权局":      77,
	"国家档案局":        79,
}

var smallDepartmentIDMap = map[string]uint{
	"国家发改委":        6,
	"国家民政部":        10,
	"国家财政部":        12,
	"国家人力资源和社会保障部": 13,
	"国家生态环境部":      15,
	"国家交通运输部":      17,
	"国家水利部":        18,
	"国家农业农村部":      19,
	"国家商务部":        20,
	"国家文化和旅游部":     21,
	"国家卫生健康委员会":    22,
	"国家退役军人事务部":    23,
	"国家海关总署":       28,
	"国家税务总局":       29,
	"国家市场监督管理总局":   30,
	"国家统计局":        34,
	"国家医疗保障局":      36,
	"国家版权局":        40,
	"中国工程院":        44,
	"国家信访局":        48,
	"国家粮食和物资储备局":   49,
	"国家能源局":        50,
	"国家移民管理局":      54,
	"国家林业和草原局":     55,
	"国家邮政局":        58,
	"国家中医药管理局":     60,
	"国家外汇管理":       64,
	"国家药监局":        65,
	"国家知识产权局":      66,
	"国家档案局":        68,
}

type ExternalSourcesMetaColly struct {
	c *colly.Collector
	// 遍历起始页
	startPages  []string
	currentPage string
	metaDal     *database.MetaDal
	dMapDal     *database.SmallDepartmentMapDal
}

func (s *ExternalSourcesMetaColly) Init() {
	s.metaDal = &database.MetaDal{Db: database.MyDb()}
	s.dMapDal = &database.SmallDepartmentMapDal{Db: database.MyDb()}
}

func (s *ExternalSourcesMetaColly) PageTraverse() {
	for key := range departmentIDMap {
		initPage := "https://bootapi.51bmj.cn/bmj-api/api/pms/pmsPolicyNoticeInfo/getDepartmentPolicyNotice?departmentId="
		requestID := departmentRequestIDMap[key]
		totalPage := getTotalPage(requestID)
		for i := 1; i <= totalPage; i++ {
			s.startPages = append(s.startPages, initPage+requestID+"&pageNo="+strconv.Itoa(i))
		}
	}
	//for _, did := range departmentRequestIDMap {
	//	initPage := "https://bootapi.51bmj.cn/bmj-api/api/pms/pmsPolicyNoticeInfo/getDepartmentPolicyNotice?departmentId="
	//	totalPage := getTotalPage(did)
	//	for i := 1; i <= totalPage; i++ {
	//		s.startPages = append(s.startPages, initPage+did+"&pageNo="+strconv.Itoa(i))
	//	}
	//}
}

func (s *ExternalSourcesMetaColly) Operate() {

	// 发送GET请求
	response, err := http.Get(s.currentPage)
	if err != nil {
		fmt.Printf("发生错误：%s\n", err)
		return
	}
	defer response.Body.Close()

	// 读取响应体
	body, err := io.ReadAll(response.Body)
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

	for _, record := range resp.Result.DepartmentPolicyNotice.Records {
		url := getUrl(record.EntranceURL)
		if url == "" {
			fmt.Println(string(debugRespBody))
			continue
		}
		dateTime, _ := utils.StringToTime(record.ReleaseTime)
		did := departmentIDMap[resp.Result.DepartmentName]
		if did == 0 {
			fmt.Println("did = 0")
		}
		metaID := s.metaDal.InsertMeta(dateTime, record.Title, url, did, provinceID)
		//用 s.metaDal.InsertMeta 方法，插入元数据，包括时间、标题、URL、部门ID和省份ID，并将返回的元数据ID存储在 metaID 中。

		sdID := smallDepartmentIDMap[resp.Result.DepartmentName]
		if sdID == 0 {
			fmt.Println("sdID = 0")
		}
		s.dMapDal.InsertDID(metaID, sdID) //调用 s.dMapDal.InsertDID 方法，插入元数据ID和小部门ID

		fmt.Printf("标题：%s，URL：%s，日期：%s\n", record.Title, url, dateTime)
	}

	//<-time.After(1 * time.Second)

}

func (s *ExternalSourcesMetaColly) Run() {
	for _, page := range s.startPages {
		s.currentPage = page
		s.Operate()
	}
}

func (s *ExternalSourcesMetaColly) Destroy() {
	// 下次运行是在一天后了，指向nil，保证内存释放，让gc自动去回收
	s.c = nil
	s.metaDal = nil
	s.startPages = nil
}

func (s *ExternalSourcesMetaColly) ExecuteWorkflow() {
	s.Init()
	s.PageTraverse()
	s.Run()
	s.Destroy()
}

func getTotalPage(id string) int {
	// 发送GET请求
	response, err := http.Get("https://bootapi.51bmj.cn/bmj-api/api/pms/pmsPolicyNoticeInfo/getDepartmentPolicyNotice?departmentId=" + id)
	if err != nil {
		fmt.Printf("发生错误：%s\n", err)
		return 0
	}
	defer response.Body.Close() //使用 defer 关键字确保在函数返回前关闭响应体。

	// 读取响应体
	body, err := io.ReadAll(response.Body) //读取响应体的全部内容，并将其存储在 body 变量中
	if err != nil {
		fmt.Printf("读取响应体时发生错误：%s\n", err)
		return 0
	}

	// 解析JSON
	var data map[string]interface{} //试将响应体的内容解析为一个 map[string]interface{} 类型的变量 data
	if err := json.Unmarshal(body, &data); err != nil {
		fmt.Printf("解析JSON时发生错误：%s, body: %v\n", err, body)
		return 0
	}

	var resp struct {
		Result struct {
			DepartmentPolicyNotice struct {
				Pages int `json:"pages"`
			} `json:"departmentPolicyNotice"`
		} `json:"result"`
	} //定义了一个嵌套的结构体 resp，用于存储解析后的JSON数据。

	// 解析JSON
	if err := json.Unmarshal(body, &resp); err != nil {
		fmt.Printf("解析JSON时发生错误：%s, body: %v\n", err, body)
		return 0
	} //再次尝试将响应体的内容解析为之前定义的结构体 resp。如果解析过程中发生错误，函数会打印错误信息并返回0。
	return resp.Result.DepartmentPolicyNotice.Pages
}

var (
	getUrlColly *colly.Collector
	result      string
)

var debugRespBody []byte

func getUrl(url string) string {
	if getUrlColly == nil {
		getUrlColly = colly.NewCollector()
		getUrlColly.Async = false

		getUrlColly.OnHTML(".from_b a[href]", func(e *colly.HTMLElement) {
			result = e.Attr("href")
			fmt.Printf("Link found: %q -> %s\n", e.Text, result)
		})
		// 记录body，用于调试问题
		getUrlColly.OnResponse(func(r *colly.Response) {
			debugRespBody = r.Body
		})
	}
	result = ""
	// 访问指定的 URL
	getUrlColly.Visit(url)

	return result
}

var _ service.MetaCrawler = (*ExternalSourcesMetaColly)(nil)
