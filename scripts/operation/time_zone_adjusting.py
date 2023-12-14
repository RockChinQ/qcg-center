import os 
import sys
import datetime

import pymongo

uri = os.environ['MONGO_URI']

client = pymongo.MongoClient(uri)

db = client['qcg-center-records']

# 遍历每一个记录，如果 begin 的小时是 0 那么将 begin 减去 8 小时

for record in db['analysis_daily'].find():
    print(record)
    begin = record['begin']
    if begin.hour == 0:
        record['begin'] = begin - datetime.timedelta(hours=8)
        result = db['analysis_daily'].update_one(
            {'_id': record['_id']},
            {'$set': {'begin': record['begin']}}
        )
        print(result)
        print(record['_id'])
        print(record['begin'])
        print('------------------')

        # input('press enter to continue...')
