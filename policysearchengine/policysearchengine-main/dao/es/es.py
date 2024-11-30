from elasticsearch import Elasticsearch

import json
import requests
from datetime import datetime
from typing import List
import yaml
import warnings

# warnings.filterwarnings("ignore")

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
	  }
	}
	'''
ResultFmtDSL = '''
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
	}'''

fmtScoreFilter = '''
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
	}'''

fmtMustDsl = '''
	{
	  "bool": {
	    "should": [
	      { "match": { "title": "%s" }},
	      { "match": { "content": "%s" }}
	    ],
	    "minimum_should_match": 1
	  }
	}'''

exactFmt = '''
	{
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
	}'''


class ESDocument:
    def __init__(self, title: str, url: str, date: datetime, content: str,
                 department_id: int, sdp_ids: List[int], province_id: int):
        self.title = title
        self.url = url
        self.date = date.isoformat()  # 将日期转换为 ISO 格式字符串
        self.content = content
        self.department_id = department_id
        self.sdp_ids = sdp_ids
        self.province_id = province_id


# def init_es_client(config_path=None):
#     """初始化并返回Elasticsearch客户端"""
#     # 初始化配置解析器
#     config = configparser.ConfigParser()
#     # 读取配置文件
#     config.read(config_path)
#     # 从配置文件中获取Elasticsearch配置
#     es_host = config.get("es", "addr")
#     # es_user = config.get('elasticsearch', 'ES_USER')
#     # es_password = config.get('elasticsearch', 'ES_PASSWORD')
#
#     es = Elasticsearch(
#         hosts=[es_host],
#         # basic_auth=(es_user, es_password),
#         verify_certs=False,
#         ca_certs='conf/http_ca.crt'
#     )
#     return es
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


def create_index(es, index_name="test-index"):
    """创建索引，如果索引已存在则忽略"""
    if not es.indices.exists(index=index_name):
        es.indices.create(index=index_name)


def define_mapping(es, index_name="test-index"):
    """为索引定义映射"""
    mapping = {
        "mappings": {
            "properties": {
                "name": {"type": "text"},
                "age": {"type": "integer"},
                "email": {"type": "keyword"}
            }
        }
    }
    es.indices.create(index=index_name, body=mapping, ignore=400)  # ignore=400忽略索引已存在错误


def insert_document(es, index_name="test-index", doc_id=None, document=None):
    """插入文档到指定索引"""
    es.index(index=index_name, id=doc_id, document=document)


def update_document(es, index_name="test-index", doc_id=None, updated_doc=None):
    """更新指定ID的文档"""
    es.update(index=index_name, id=doc_id, body={"doc": updated_doc})


def delete_document(es, index_name="test-index", doc_id=None):
    """删除指定ID的文档"""
    es.delete(index=index_name, id=doc_id)


def search_documents(es, index_name="test-index", query=None):
    """在指定索引中搜索文档"""
    return es.search(index=index_name, body=query)


def index_doc(date: datetime, department_id: int, province_id: int,
              title: str, url: str, content: str, sdp_ids: List[int], index: str):
    doc = ESDocument(title, url, date, content, department_id, sdp_ids, province_id)

    data = json.dumps(doc.__dict__)  # 将文档对象转换为 JSON 格式
    # index = 'your_index_name'  # 这里替换为你的目标索引名称
    es_url = f'http://localhost:9200/{index}/_doc'  # Elasticsearch 文档上传的 URL

    response = requests.post(es_url, data=data, headers={"Content-Type": "application/json"})  # 上传数据

    if response.status_code != 201:
        print(f"ES上传数据失败，err: {response.text}")
        return

    print(response.json())  # 打印索引响应的内容
