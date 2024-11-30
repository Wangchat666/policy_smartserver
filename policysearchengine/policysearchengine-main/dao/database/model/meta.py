from sqlalchemy import Column, Integer, DateTime, String, ForeignKey
from sqlalchemy.orm import declarative_base
import datetime

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
