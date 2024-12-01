package model

import "time"

type ESDocument struct {
	Title             string    `json:"title"`
	Url               string    `json:"url"`
	Date              time.Time `json:"date"`
	Content           string    `json:"content"`
	DepartmentID      uint      `json:"department_id"`
	SmallDepartmentID []uint    `json:"small_department_id"`
	ProvinceID        uint      `json:"province_id"`
}

type ESResp struct {
	Hits struct {
		Hits []struct { //Hits字段是一个切片，包含了一系列的搜索结果。
			//ID     string      `json:"_id"`
			//Index  string      `json:"_index"`
			//Score  interface{} `json:"_score"`
			Source struct {
				//Content      string    `json:"content"`
				Date              time.Time `json:"date"`
				DepartmentID      int       `json:"department_id"`
				ProvinceID        int       `json:"province_id"`
				SmallDepartmentID []int     `json:"small_department_id"`
				Title             string    `json:"title"`
				URL               string    `json:"url"`
			} `json:"_source"` //结束Source结构体的定义，并指定JSON反序列化时使用_source作为键。

			Highlight struct {
				Title   []string `json:"title"`
				Content []string `json:"content"`
			} `json:"highlight"`
		} `json:"hits"`
		Total struct {
			Relation string `json:"relation"`
			Value    int    `json:"value"`
		} `json:"total"`
	} `json:"hits"`
	TimedOut bool `json:"timed_out"`
	Took     int  `json:"took"`
} //解析Elasticsearch搜索响应的JSON数据。这个结构体包含了Elasticsearch响应的不同字段，使得可以方便地反序列化JSON响应，并访问其中的数据。
