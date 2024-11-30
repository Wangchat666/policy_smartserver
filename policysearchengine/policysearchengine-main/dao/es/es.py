from elasticsearch import Elasticsearch
import json
import requests
from datetime import datetime
from typing import List
import yaml
# import warnings

# warnings.filterwarnings("ignore")
index = 'policy_search'
query_fmt = '''
	{
	  "query": {
	    "bool": {
	      "must": [
	        {
	          "bool": {
	            "should": [
	              {
	                "match": {
	                  "title": {}
	                }
	              },
	              {
	                "match": {
	                  "content": {}
	                }
	              }
	            ]
	          }
	        },
	        { "bool": {"must": [{}]} }
	        {}
	      ]
	    }
	  }
	}
	'''
ResultFmtDSL = '''
	{
	  {}
	  "sort": [
	    {
	      "date": {
	        "order": "asc"
	      }
	    }
	  ],
	  "from": {},
	  "size": {},
	  "highlight": {}
	}
	'''

# 构建高亮查询
highlight = '''
	{
	  "fields": {
	    "title": {},
	    "content": {}
	  },
	  "fragment_size": 50,
	  "pre_tags": ["<em style='color:red'>"],
	  "post_tags": ["</em>"]
	}
	'''

# 分数搜索DSL
fmtScoreQuery = '''
	{
	  "query": {
	    "function_score": {
	      {}
	      "functions": [
	        {}
	      ],
	      "score_mode": "multiply",
	      "boost_mode": "multiply"
	    }
	  },
	  "from": {},
	  "size": {},
	  "highlight": {}
	}'''

fmtScoreFilter = '''
	{
	  "filter": {
	    "bool": {
	      "should": [
	        { "match": { "title": {} }},
	        { "match": { "content": {} }}
	      ]
	    }
	  },
	  "weight": {}
	}'''

fmtMustDsl = '''
	{
	  "bool": {
	    "should": [
	      { "match": { "title": {} }},
	      { "match": { "content": {} }}
	    ],
	    "minimum_should_match": 1
	  }
	}'''

exactFmt = '''
	{
	  "bool": {
	    "should": [{
	      "match_phrase": {
	        "title": {}
	      }
	    }, {
	      "match_phrase": {
	        "content": {}
	      }
	    }]
	  }
	}'''

class SearchInput:
    def __init__(self, text: str, use_score: bool, score_field: dict):
        self.text = text  # 搜索内容
        self.use_score = use_score
        self.score_field = score_field  # 权重分数
class ESDocument:
    def __init__(self, title: str, labels:str,url: str, date: datetime, content: str,abstract: str,keywords: str,classification: str,
                 department_id: int, sdp_ids: List[int], province_id: int):
        self.title = title
        self.labels = labels
        self.url = url
        self.date = date.isoformat()  # 将日期转换为 ISO 格式字符串
        self.content = content
        self.abstract = abstract
        self.keywords = keywords
        self.classification = classification
        self.department_id = department_id
        self.sdp_ids = sdp_ids
        self.province_id = province_id



def init(config_path:str):
    with open(config_path, 'r',encoding='utf-8') as file:
        config = yaml.safe_load(file)  # 读取 YAML 文件并解析为字典

    es_address = config['es']['addr']  # 获取 Elasticsearch 地址
    # es_index = config['es']['index']  # 获取索引名称

    es = Elasticsearch([es_address])  # 创建 Elasticsearch 客户端

    if not es.ping():
        print(f"ES初始化错误，无法连接到 {es_address}")
        return 0

    print(f"ES客户端已初始化，连接到 {es_address}")
    return es  # 返回 ES 客户端


def index_doc(title: str,labels:str, url: str,date: datetime, content: str,department_id: int,
              abstract: str,keywords: str,classification: str,sdp_ids: List[int],province_id: int, index: str):
    doc = ESDocument(title, labels,url, date, content,  abstract, keywords, classification,department_id, sdp_ids, province_id)

    data = json.dumps(doc.__dict__)  # 将文档对象转换为 JSON 格式
    # index = 'your_index_name'  # 这里替换为你的目标索引名称
    es_url = f'http://localhost:9200/{index}/_doc'  # Elasticsearch 文档上传的 URL

    response = requests.post(es_url, data=data, headers={"Content-Type": "application/json"})  # 上传数据

    if response.status_code != 201:
        print(f"ES上传数据失败，err: {response.text}")
        return

    print(response.json())  # 打印索引响应的内容



def match_all_doc(es:Elasticsearch,index: str):
    query = {
        "query": {
            "match_all": {}
        }
    }

    # 执行搜索
    search = es.search(index=index, body=query)
    return search

def sprintf(format_str: str, *args) -> str:
    """
    使用给定的格式字符串和参数生成格式化后的字符串。

    :param format_str: 格式字符串，类似于 Python 中的字符串格式化
    :param args: 可变数量的参数，用于格式化
    :return: 生成的格式化字符串
    """
    return format_str.format(*args)


def query_fmt_print(text, exact, small_department_id, province_id):
    query = ","

    if small_department_id == 0 and province_id == 0:
        query = ""
    elif small_department_id != 0 and province_id != 0:
        query += '{{ "match": {{ "small_department_id": {} }} }},'.format(small_department_id)
        query += '{{ "match": {{ "province_id": {} }} }}'.format(province_id)
    elif small_department_id != 0 and province_id == 0:
        query += '{{ "match": {{ "small_department_id": {} }} }}'.format(small_department_id)
    elif small_department_id == 0 and province_id != 0:
        query += '{{ "match": {{ "province_id": {} }} }}'.format(province_id)

    exact_trim = [s for s in exact if s != ""]  # 去除空字符串，得到 exact_trim 列表

    exact_query = ""
    for i in range(len(exact_trim)):
        exact_query += sprintf(exactFmt, exact_trim[i])
        if i < len(exact_trim) - 1:
            exact_query += ","
    return sprintf(query_fmt, text, text, exact_query, query)

def SearchDocWithSmallDepartmentID(es:Elasticsearch,searchquery: SearchInput, exact:List[str], smalldepartmentID:int, provinceID:int, ffrom:int, size:int):
    searchFmt = query_fmt_print(searchquery.text, exact, smalldepartmentID, provinceID)

    if not searchquery.use_score:  # 如果UseScore为false
        query = sprintf(ResultFmtDSL,
			searchFmt,
			ffrom-1,
			size,
			highlight)
    else:
        query = sprintf(fmtScoreQuery,
			searchFmt,
			fmtScoreFilters(searchquery.score_field),
			ffrom-1,
			size,
			highlight)  # 假设fmt_score_query和fmt_score_filters已定义

    print(query)
    query=json.loads(query)
    # 发送Elasticsearch搜索请求
    search_result = es.search(
        index=index,  # 假设index已定义
        body=query,  # 将查询转换为JSON格式
        track_total_hits=True  # 启用跟踪总命中数的选项
    )
    if search_result.get('error'):
        print("Error executing search:", search_result['error'])
        return None

    # 解析search_result中的JSON数据
    response_data = json.loads(search_result['body'])  # search_result['body']是一个JSON字符串

    return response_data

def fmtScoreFilters(m):
    filters = []
    for word, weight in m.items():
        filter_string = sprintf(fmtScoreFilter,word, word, weight)  # 假设fmt_score_filter已定义
        filters.append(filter_string)
    return ", ".join(filters)  # 将所有过滤器字符串用逗号和空格连接









# es = init('.../config/config.yaml')

