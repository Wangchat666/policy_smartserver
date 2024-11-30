# -*- coding: utf-8 -*-
import pymysql
import yaml
with open('../../config/config.yaml', 'r', encoding='utf-8') as f:
    V = yaml.load(f, Loader=yaml.FullLoader)


# 创建数据库连接
db = pymysql.connect(
    host='localhost',  # MySQL服务器地址
    user=V['mysql']['user'],  # 用户名
    password=V['mysql']['password'],  # 密码
    database=V['mysql']['dbname']  # 数据库名称
)

# 创建游标对象
cursor = db.cursor()
# sql = "CREATE DATABASE policy_search "
# cursor.execute(sql)