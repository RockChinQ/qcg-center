# 从installer-reports中清洗数据
# - 过时的数据中，有些文档中的remote_addr是包含端口号的，删除端口号

import os
import pymongo

uri = os.environ['MONGO_URI']

client = pymongo.MongoClient(uri)

# 从installer-reports中清洗数据
# 查找remote_addr中包含:的文档
# 挨个执行：删除:后面的部分并更新文档
all_legacy_docs = client['qcg-center-records']['installer-reports'].find({'remote_addr': {'$regex': r':\d+$'}})

for doc in all_legacy_docs:
    
    remote_addr = doc['remote_addr']
    remote_addr = remote_addr.split(':')[0]
    doc['remote_addr'] = remote_addr
    client['qcg-center-records']['installer-reports'].update_one({'_id': doc['_id']}, {'$set': {'remote_addr': remote_addr}})
    
    print("Updated", remote_addr)
    
# print("Found", len(to_delete), "documents to delete")
# 删除installer-reports中timestamp小于2022-12-01的文档
import datetime
ts = datetime.datetime(2022, 12, 1).timestamp()

client['qcg-center-records']['installer-reports'].delete_many({'timestamp': {'$lt': ts}})

