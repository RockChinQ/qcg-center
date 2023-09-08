import os
import pymongo

uri = os.environ['MONGO_URI']

# 删除qcg-center-records.qchatgpt-usage中所有source为migrationd的文档

client = pymongo.MongoClient(uri)

db = client['qcg-center-records']

db['qchatgpt-usage'].delete_many({'source': 'migration2'})