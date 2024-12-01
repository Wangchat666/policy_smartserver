package http

import (
	"PolicySearchEngine/config"
	"PolicySearchEngine/dao/database"
	"PolicySearchEngine/dao/es"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const (
	size         = 10
	SummaryFault = "大语言模型异常，分析摘要失败..."
)

func Router() {
	router := gin.Default()
	router.GET("/search", searchHandel)
	router.GET("/summary", summaryHandel)
	if err := router.Run(":8080"); err != nil {
		fmt.Printf("router.Run Failed, err:%+v\n", err)
	}
}

func summaryHandel(c *gin.Context) {
	// 根据URL查询到对应文章内容
	u := c.Query("url")
	metaDal := &database.MetaDal{Db: database.MyDb()}
	meta := metaDal.GetMetaByUrl(u)
	contentDal := &database.ContentDal{Db: database.MyDb()}
	content := contentDal.GetContentByMetaID(meta.ID)

	// 调用大语言模型接口处理文章内容
	result := docSummary(content.Article)
	fmt.Println(result)
	c.JSON(http.StatusOK, gin.H{
		"summary": result,
	})
}

func searchHandel(c *gin.Context) {

	searchQuery := c.Query("s")
	exactQuery := c.Query("e")
	sDepartmentID := c.Query("did")
	sProvinceID := c.Query("pid")
	sPage := c.Query("page")

	exact := strings.Split(exactQuery, ",")
	page, err1 := strconv.Atoi(sPage)
	departmentID, err2 := strconv.Atoi(sDepartmentID)
	provinceID, err3 := strconv.Atoi(sProvinceID)
	if err1 != nil || err2 != nil || err3 != nil {
		fmt.Printf("Get请求参数错误\n，err1:%+v，err2:%+v，err3:%+v", err1, err2, err3)
		return
	}

	// 检查查询参数是否为空
	if searchQuery == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少search参数"})
		return
	}

	//esResp := es.SearchDocByDepartmentID(preSearch(searchQuery), departmentID, provinceID, page, size)
	esResp := es.SearchDocWithSmallDepartmentID(preSearch(searchQuery), exact, departmentID, provinceID, page, size)
	//totalPage := math.Ceil(float64(esResp.Hits.Total.Value) / size)

	c.JSON(http.StatusOK, esResp.Hits)
}

func preSearch(s string) es.SearchInput {
	result := es.SearchInput{
		Text:     s,
		UseScore: false,
	}
	preSearchUrl := config.V.GetString("http.pre-search-url") + url.QueryEscape(s)
	resp, err := http.Get(preSearchUrl)
	if err != nil {
		fmt.Println("Error occurred while sending GET request:", err)
		return result
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return result
	}

	var r struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}

	err = json.Unmarshal(body, &r)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return result
	}

	if resp.StatusCode == http.StatusOK && r.Code == 200 {
		scoreField, ok := parseKeywords(r.Message)
		if ok {
			var fields []string
			for field := range scoreField {
				fields = append(fields, field)
			}

			return es.SearchInput{
				Text:       strings.Join(fields, " "),
				ScoreField: scoreField,
				UseScore:   true,
			}
		} else {
			return result
		}
	} else {
		fmt.Printf("Non-200 status code: %d\n, resp:%+v", resp.StatusCode, r)
	}
	return result
}

func docSummary(s string) string {

	payload, err := json.Marshal(map[string]string{
		"message": s,
	})
	if err != nil {
		fmt.Println("JSON编码失败:", err)
		return SummaryFault
	}

	resp, err := http.Post(config.V.GetString("http.doc-summary-url"),
		"application/json",
		bytes.NewBuffer(payload),
	)
	if err != nil {
		fmt.Println("POST请求失败:", err)
		return SummaryFault
	}
	defer resp.Body.Close()

	// 读取并解析响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("读取响应失败:", err)
		return SummaryFault
	}

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		fmt.Printf("JSON decoding failed, body:%+v, err:%+v\n", string(body), err)
		return SummaryFault
	}
	return data["message"].(string)
}

func parseKeywords(text string) (map[string]float64, bool) {
	keywords := make(map[string]float64)

	parts := strings.Split(text, ",")
	if len(parts) == 0 {
		return nil, false
	}
	for _, part := range parts {
		kv := strings.Split(part, ":")
		if len(kv) != 2 {
			return nil, false
		}
		keyword := strings.TrimSpace(kv[0])
		weightStr := strings.TrimSpace(kv[1])
		weight, err := strconv.ParseFloat(weightStr, 64)
		if err != nil {
			return nil, false
		}
		keywords[keyword] = weight
	}

	return keywords, true
}
