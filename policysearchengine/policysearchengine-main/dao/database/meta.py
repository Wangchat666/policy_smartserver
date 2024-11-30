from sqlalchemy.orm import Session, sessionmaker
from model.meta import Meta
import init
import logging


class MetaDal:
    def __init__(self, db: Session):
        self.db = db

    def insert_meta(self, date, title: str, url: str, department_id: int, province_id: int) -> int:
        meta = Meta(
            date=date,
            title=title,
            url=url,
            department_id=department_id,
            province_id=province_id,
        )

        # 使用 SQLAlchemy 的方法确保唯一性（根据 URL）
        existing_meta = self.db.query(Meta).filter(Meta.url == url).first()
        if existing_meta:
            return existing_meta.id

        try:
            self.db.add(meta)
            self.db.commit()
            self.db.refresh(meta)
            return meta.id
        except Exception as e:
            logging.error(f"InsertMeta... {date}, {meta}, Error: {e}")
            self.db.rollback()
            return None

    def update_meta_title(self, title: str, url: str):
        try:
            meta = self.db.query(Meta).filter(Meta.url == url).first()
            if meta is not None:
                meta.title = title
                self.db.commit()
            else:
                logging.warning(f"Meta with URL {url} not found for update.")
        except Exception as e:
            logging.error(f"UpdateMetaTitle... {title}, {url}, Error: {e}")
            self.db.rollback()

    def get_all_meta(self, department_id: int, province_id: int) -> list:
        try:
            return self.db.query(Meta).filter(Meta.department_id == department_id,
                                              Meta.province_id == province_id).all()
        except Exception as e:
            logging.error(f"读取数据失败: {e}")
            return []

    def get_all_meta_by_ids(self, province_id: int, id: int) -> list:
        try:
            return self.db.query(Meta).filter(Meta.province_id == province_id, Meta.id > id).all()
        except Exception as e:
            logging.error(f"读取数据失败: {e}")
            return []

    def get_meta_by_url(self, url: str) -> Meta:
        try:
            return self.db.query(Meta).filter(Meta.url == url).first()
        except Exception as e:
            logging.error(f"读取数据失败: {e}")
            return None


Session = sessionmaker(bind=init.db)
db = Session()
M = MetaDal(db)
M.insert_meta("2024-03-13", "我国支持科技创新主要税费优惠政策指引",
              "https://www.most.gov.cn/kjbgz/202403/t20240313_189961.html", 1, 35)
