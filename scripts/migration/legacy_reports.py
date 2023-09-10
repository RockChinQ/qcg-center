import os
import sys

MYSQL_HOST = os.environ.get('MYSQL_HOST', 'localhost')
MYSQL_PORT = os.environ.get('MYSQL_PORT', '3306')
MYSQL_USER = os.environ.get('MYSQL_USER', 'root')
MYSQL_PASSWORD = os.environ.get('MYSQL_PASSWORD', 'password')
MYSQL_DB = os.environ.get('MYSQL_DB', 'test')

uri = os.environ['MONGO_URI']

import pymongo

client = pymongo.MongoClient(uri)

import pymysql

conn = pymysql.connect(host=MYSQL_HOST, port=int(MYSQL_PORT), user=MYSQL_USER, passwd=MYSQL_PASSWORD, db=MYSQL_DB)

print("Connected to MySQL")

cursor = conn.cursor()

offset = os.environ.get('OFFSET', 0)

limit = 4000

# 数据量很大，一条一条处理
# 不要fetchall，会占用大量内存
try:
    while True:
        cursor.execute(
            "select `id`, `osname`, `arch`, `timestamp`, `mac`, `version`, `message`, `ip` as `remote_addr` from `reports` order by id asc limit %s offset %s",
            (limit,offset,))
        
        # 判空
        if cursor.rowcount == 0:
            break
        
        to_insert = []
        
        id = 0
        
        for row in cursor.fetchall():
            id, osname, arch, timestamp, mac, version, message, remote_addr = row
            
            to_insert.append({
                'osname': osname,
                'arch': arch,
                'timestamp': int(timestamp),
                'mac': mac,
                'version': version,
                'message': message,
                'remote_addr': remote_addr,
            })
            
        if len(to_insert) == 0:
            break
        
        client['qcg-center-records']['installer-reports'].insert_many(to_insert)
        
        print("Processed offset", offset)
        
        offset += limit
        
except Exception as e:
    import traceback
    traceback.print_exc()
    
    print("Error at offset", offset)