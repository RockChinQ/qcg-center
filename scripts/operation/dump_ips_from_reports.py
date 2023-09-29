import datetime
import os
import pymongo

uri = os.environ['MONGO_URI']

client = pymongo.MongoClient(uri)

db = client['qcg-center-records']

result = db['installer-reports'].aggregate(
    [
        {
            "$group": {
                "_id": "$remote_addr",
                "first_doc": {"$first": "$$ROOT"}
            }
        },
        {
            "$replaceRoot": {
                "newRoot": "$first_doc"
            }
        },
        {
            "$project": {
                "_id": 0,
                "remote_addr": 1,
                "timestamp": 1,
            }
        },
        {
            "$sort": {
                "timestamp": 1
            }
        }
    ]
)

result = list(result)

# print(result)

unique_remote_addrs = [res['remote_addr'] for res in result]

print(len(unique_remote_addrs))

to_insert = []
"""
{
    "remote_addr": "xxx"
    "created_at": 从timestamp(秒)创建时间类型,
}
"""

for res in result:
    to_insert.append({
        "remote_addr": res['remote_addr'],
        "created_at": datetime.datetime.fromtimestamp(res['timestamp'])
    })
    
db['analysis_reports_remote_addrs'].insert_many(to_insert)