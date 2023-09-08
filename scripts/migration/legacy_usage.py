import os
import sys

MYSQL_HOST = os.environ.get('MYSQL_HOST', 'localhost')
MYSQL_PORT = os.environ.get('MYSQL_PORT', '3306')
MYSQL_USER = os.environ.get('MYSQL_USER', 'root')
MYSQL_PASSWORD = os.environ.get('MYSQL_PASSWORD', 'password')
MYSQL_DB = os.environ.get('MYSQL_DB', 'test')

config_file = sys.argv[1]

import yaml

cfg = yaml.load(open(config_file), Loader=yaml.FullLoader)

import pymongo

uri = cfg['database']['params']['uri']

client = pymongo.MongoClient(uri)

import pymysql

conn = pymysql.connect(host=MYSQL_HOST, port=int(MYSQL_PORT), user=MYSQL_USER, passwd=MYSQL_PASSWORD, db=MYSQL_DB)

print("Connected to MySQL")

cursor = conn.cursor()

# cursor.execute('SELECT `id`,`service_name`,`version`,`count`,`msg_source`,`timestamp`,`ip` as `remote_addr` FROM `usage` order by id asc')

print("MySQL query done")

# 检查是否有进度
offset = 0

if os.path.exists('legacy_usage.log'):
    with open('legacy_usage.log', 'r') as f:
        offset = int(f.read())
        print("Found offset", offset)

limit = 4000

# 数据量很大，一条一条处理
# 不要fetchall，会占用大量内存
try:
    while True:
        cursor.execute(
            'SELECT `id`,`service_name`,`version`,`count`,`msg_source`,`timestamp`,`ip` as `remote_addr` FROM `usage` order by id asc limit %s offset %s',
            (limit,offset,))
        
        # 判空
        if cursor.rowcount == 0:
            break
        
        to_insert = []
        
        id = 0
        
        for row in cursor.fetchall():
            id, service_name, version, count, msg_source, timestamp, remote_addr = row
            to_insert.append({
                'service_name': service_name,
                'version': version,
                'count': int(count),
                'msg_source': msg_source,
                'timestamp': int(timestamp),
                'remote_addr': remote_addr,
                'source': "migration2"
            })
        
        # id_max, service_name, version, count, msg_source, timestamp, remote_addr = cursor.fetchone()
        # print(id_max, service_name, version, count, msg_source, timestamp, remote_addr)
        
        # client['qcg-center-records']['qchatgpt-usage'].insert_one({
        #     'service_name': service_name,
        #     'version': version,
        #     'count': int(count),
        #     'msg_source': msg_source,
        #     'timestamp': int(timestamp),
        #     'remote_addr': remote_addr,
        #     'source': "migration2"
        # })
        
        id_max = id
        
        client['qcg-center-records']['qchatgpt-usage'].insert_many(to_insert)
        
        print(id_max, "inserted")
        
        offset += limit
except Exception as e:
    import traceback
    
    traceback.print_exc()
    # 保存进度
    with open('legacy_usage.log', 'w') as f:
        f.write(str(offset))
