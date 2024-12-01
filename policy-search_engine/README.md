# PolicySearchEngine 

政汇搜索引擎 ——旨在开发一个全面爬取政策文件并提供搜索引擎服务的工具。

目前爬虫与搜索引擎部分已经完成，精确检索与前端web页面正在完善中。

开发进度见：[开发日志](doc/开发日志.md)

## QuickStart

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
