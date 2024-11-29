# PolicySearchEngine 

毕设题目：面向公共政策文本检索的搜索引擎设计与实现

目前写的是爬虫部分，等爬虫写好后把整个项目结构改改，然后再做web部分

开发进度见：[开发日志](doc/开发日志.md)

## 运行

```shell
# 前端
cd ./front-app
npm start

# 大语言模型接口
cd ./pre-search
python main.py

# 启动es
D:\download\elasticsearch-8.12.0\bin\elasticsearch.bat

# 爬虫&后端接口
go run main.go
```
