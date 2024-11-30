
from dao.es.es import init


# def main():
#     ES = init('./config/config.yaml')
if __name__ == '__main__':
    config_path = './config/config.yaml'  # 配置文件的路径
    es_client = init(config_path)  # 调用 init 函数并获取客户端实例

