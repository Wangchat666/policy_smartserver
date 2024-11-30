from datetime import datetime
from dao.es.es import init,index_doc


# def main():
#     ES = init('./config/config.yaml')
if __name__ == '__main__':
    config_path = './config/config.yaml'  # 配置文件的路径
    es_client = init(config_path)# 调用 init 函数并获取客户端实例
    date_instance = datetime.strptime("2024.11.30", "%Y.%m.%d")
    index_doc( 'policy_smartserver', '自我介绍',
              'https://blog.csdn.net/weixin_46264660/article/details/130238426',date_instance,"我正帅",
               7,"帅","帅","帅",[1,2,3,4,5,6,7],7,"test")

