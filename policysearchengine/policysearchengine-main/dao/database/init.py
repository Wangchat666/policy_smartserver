from sqlalchemy import create_engine
from sqlalchemy.exc import IntegrityError
from sqlalchemy import Column, Integer, String, DateTime, ForeignKey, Text, Table
from sqlalchemy.orm import declarative_base
from urllib.parse import quote
import datetime
import logging
import yaml

with open('../../config/config.yaml', 'r', encoding='utf-8') as f:
    V = yaml.load(f, Loader=yaml.FullLoader)
# 设置日志记录
logging.basicConfig(level=logging.ERROR)

Base = declarative_base()


class Meta(Base):
    __tablename__ = 'meta'

    id = Column(Integer, primary_key=True)  # 主键列
    created_at = Column(DateTime, default=datetime.datetime.utcnow)  # 创建时间
    updated_at = Column(DateTime, index=True, default=datetime.datetime.utcnow,
                        onupdate=datetime.datetime.utcnow)  # 更新时间
    date = Column(DateTime)  # 日期
    title = Column(String(255))  # 标题
    url = Column(String(255), unique=True)  # 唯一 URL
    labels = Column(String(255))  # 标签
    classification = Column(String(255))  # 分类
    department_id = Column(Integer, ForeignKey('department.id'), index=True)  # 部门ID, 假设有一个 department 表
    province_id = Column(Integer, ForeignKey('province.id'), index=True)  # 省份ID, 假设有一个 province 表


class Content(Base):
    __tablename__ = 'content'

    id = Column(Integer, primary_key=True)  # 主键列
    meta_id = Column(Integer, ForeignKey('meta.id'), nullable=False, unique=True)  # 外键，MetaID（假设 Meta 表存在）
    article = Column(Text)  # 设置为 mediumtext 的等价类型
    keyword = Column(String(50))
    abstract = Column(Text)  # 摘要
    created_at = Column(DateTime, default=datetime.datetime.utcnow)  # 创建时间
    updated_at = Column(DateTime, index=True, default=datetime.datetime.utcnow, onupdate=datetime.datetime.utcnow)


class Department(Base):
    __tablename__ = 'department'

    id = Column(Integer, primary_key=True)  # 主键列
    name = Column(String(255))  # 部门名称


class Province(Base):
    __tablename__ = 'province'

    id = Column(Integer, primary_key=True)  # 主键列
    name = Column(String(50))  # 其他字段示例


class SmallDepartmentMap(Base):
    __tablename__ = 'small_department_map'

    id = Column(Integer, primary_key=True)  # 主键列
    meta_id = Column(Integer, ForeignKey('meta.id'), nullable=False)  # 外键
    small_department_id = Column(Integer, ForeignKey('small_department.id'), nullable=False)  # 对应SmallDepartment表id


class SmallDepartment(Base):
    __tablename__ = 'small_department'

    id = Column(Integer, primary_key=True)  # 主键列
    name = Column(String(50))  # 其他字段示例


def init_table():
    # 初始化数据表
    try:
        Base.metadata.create_all(my_db)
        print("所有表创建成功！")
    except IntegrityError as e:
        logging.error(f"Error initializing tables: {e}")


def init():
    global my_db
    user = V['mysql']['user']
    passwd = V['mysql']['password']
    addr = V['mysql']['addr']
    dbname = V['mysql']['dbname']
    psd = quote(passwd)
    # 数据库配置
    connection_string = f"mysql+pymysql://{user}:{psd}@{addr}/{dbname}"

    # 连接数据库
    my_db = create_engine(connection_string, echo=False)

    init_table()  # 初始化数据表

    return my_db


db = init()
