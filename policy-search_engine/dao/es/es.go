package es

import (
	"PolicySearchEngine/config"
	"PolicySearchEngine/model"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"strings"
	"time"
)

type SearchInput struct {
	Text       string //搜索内容？
	UseScore   bool
	ScoreField map[string]float64 //权重分数？
}

var (
	es    *elasticsearch.Client //es 将被用来持有一个 Elasticsearch 客户端的实例
	index string                //用来存储一个索引的名称。在 Elasticsearch 中，索引是存储相关数据的地方，类似于关系数据库中的表。
)

const (
	queryFmt = `
  "query": {
    "bool": {
      "must": [
        {
          "bool": {
            "should": [
              {
                "match": {
                  "title": "%s"
                }
              },
              {
                "match": {
                  "content": "%s"
                }
              }
            ]
          }
        },
		{ "bool": {"must": [%s]} }
		 %s
      ]
    }
  },
`
	ResultFmtDSL = `
	{
  %s
  "sort": [
    {
      "date": {
        "order": "asc"
      }
    }
  ],
  "from": %d,
  "size": %d,
  "highlight": %s
}
`

	// 构建高亮查询
	highlight = `
	{
		"fields": {
			"title": {},
			"content": {}
		},
		"fragment_size": 50,
		"pre_tags": ["<em style='color:red'>"],
		"post_tags": ["</em>"]
	}`

	// 分数搜索DSL
	fmtScoreQuery = `
	{
	  "query": {
	    "function_score": {
	      %s
	      "functions": [
			%s
	      ],
	      "score_mode": "multiply",
	      "boost_mode": "multiply"
	    }
	  },
		"from": %d,
		"size": %d,
		"highlight": %s
	}`

	fmtScoreFilter = `
	{
	  "filter": {
	    "bool": {
	      "should": [
	        { "match": { "title": "%s" }},
	        { "match": { "content": "%s" }}
	      ]
	    }
	  },
	  "weight": %f
	}`

	fmtMustDsl = `
            {
              "bool": {
                "should": [
                  { "match": { "title": "%s" }},
                  { "match": { "content": "%s" }}
                ],
                "minimum_should_match": 1
              }
            }`

	exactFmt = `{
					"bool": {
						"should": [{
							"match_phrase": {
								"title": "%s"
							}
						}, {
							"match_phrase": {
								"content": "%s"
							}
						}]
					}
				}`
)

func Init() {
	var err error
	cfg := elasticsearch.Config{
		Addresses: []string{
			config.V.GetString("es.addr"), //Addresses: []string{"http://localhost:9200"}
		},
	}
	es, err = elasticsearch.NewClient(cfg)
	if err != nil {
		fmt.Printf("ES初始化错误，err:%+v\n", err.Error())
	}
	index = config.V.GetString("es.index")
} //貌似是创建客户端

func IndexDoc(date time.Time, departmentID, provinceID uint, title, url, content string, sdpIDs []uint) {
	doc := model.ESDocument{
		Title:             title,
		Url:               url,
		Date:              date,
		Content:           content,
		DepartmentID:      departmentID,
		SmallDepartmentID: sdpIDs,
		ProvinceID:        provinceID,
	} //这是一个结构体，用于表示Elasticsearch文档。
	data, _ := json.Marshal(doc)                       //将doc结构体序列化为JSON格式，结果存储在data变量中
	idx, err := es.Index(index, bytes.NewReader(data)) //使用 Elasticsearch 客户端的 Index 方法将 JSON 数据上传到指定的索引中。index 是目标索引的名称，bytes.NewReader(data) 将字节数组转换为可读的 io.Reader 类型。
	if err != nil {
		fmt.Printf("ES上传数据失败，err:%+v\n", err)
		return
	}
	fmt.Println(idx.String()) //打印索引响应的内容。
} //上传es数据

func MatchAllDoc() {
	query := `{ "query": { "match_all": {} } }`
	search, err := es.Search(
		es.Search.WithIndex(index),
		es.Search.WithBody(strings.NewReader(query)),
	)
	if err != nil {
		return
	}
	fmt.Println(search.String())
}

// SearchDocWithSmallDepartmentID 根据小部门筛
func SearchDocWithSmallDepartmentID(searchQuery SearchInput, exact []string, smallDepartmentID, provinceID, from, size int) *model.ESResp {
	var query string
	searchFmt := queryFmtPrint(searchQuery.Text, exact, smallDepartmentID, provinceID)
	if !searchQuery.UseScore { //如果UseScore为false，则使用fmt.Sprintf函数和ResultFmtDSL格式化字符串来构建查询
		query = fmt.Sprintf(ResultFmtDSL,
			searchFmt,
			from-1,
			size,
			highlight)
	} else {
		query = fmt.Sprintf(fmtScoreQuery,
			searchFmt,
			fmtScoreFilters(searchQuery.ScoreField),
			from-1,
			size,
			highlight)
	}

	fmt.Println(query)

	searchResult, err := es.Search(
		es.Search.WithIndex(index),                   //指定要查询的Elasticsearch索引的名称。
		es.Search.WithBody(strings.NewReader(query)), //指定查询的请求体。
		es.Search.WithTrackTotalHits(true),           //启用跟踪总命中数的选项。
	)
	if err != nil {
		fmt.Println("Error executing search:", err)
		return nil
	}

	// 解析 searchResult 中的 JSON 数据
	var responseData model.ESResp
	_ = json.NewDecoder(searchResult.Body).Decode(&responseData) //使用json.NewDecoder解码器来解析searchResult.Body中的JSON数据，并将其存储在responseData变量中。忽略任何解码错误。

	return &responseData
}

func queryFmtPrint(text string, exact []string, smallDepartmentID, provinceID int) string {
	query := ","
	if smallDepartmentID == 0 && provinceID == 0 {
		query = ""
	} else if smallDepartmentID != 0 && provinceID != 0 {
		query += fmt.Sprintf(`{ "match": { "small_department_id": %d }},`, smallDepartmentID)
		query += fmt.Sprintf(`{ "match": { "province_id": %d }}`, provinceID)
	} else if smallDepartmentID != 0 && provinceID == 0 {
		query += fmt.Sprintf(`{ "match": { "small_department_id": %d }}`, smallDepartmentID)
	} else if smallDepartmentID == 0 && provinceID != 0 {
		query += fmt.Sprintf(`{ "match": { "province_id": %d }}`, provinceID)
	}

	var exactTrim []string
	for _, str := range exact {
		if str != "" {
			exactTrim = append(exactTrim, str)
		} //去除空字符串，得到 exactTrim 列表。
	}

	exactQuery := ""
	for i := 0; i < len(exactTrim); i++ {
		exactQuery += fmt.Sprintf(exactFmt, exactTrim[i])
		if i < len(exactTrim)-1 {
			exactQuery += ","
		}
	} //遍历 exactTrim 列表，构建精确匹配的查询字符串，并在每个条件之间添加逗号分隔。

	return fmt.Sprintf(queryFmt, text, text, exactQuery, query)
} //构建一个用于 Elasticsearch 查询的格式化字符串。通过不同的条件判断和字符串拼接，生成最终的查询字符串。

func fmtScoreFilters(m map[string]float64) string {
	var filters []string
	for word, weight := range m {
		filter := fmt.Sprintf(fmtScoreFilter, word, word, weight)
		filters = append(filters, filter)
	}
	return strings.Join(filters, ", ")
}
