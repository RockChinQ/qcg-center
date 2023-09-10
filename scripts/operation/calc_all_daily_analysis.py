import os
import datetime

import pymongo

uri = os.environ['MONGO_URI']

client = pymongo.MongoClient(uri)

db = client['qcg-center-records']

colletion = db['qchatgpt-usage']

print("Connected to MongoDB")
# 
# [
#   {
#     $addFields: {
#       cst_ts: {
#         $add: ["$timestamp", 8 * 3600],
#       },
#     },
#   },
#   {
#     $addFields: {
#       cst_date: {
#         $toDate: {
#           $multiply: ["$cst_ts", 1000],
#         },
#       },
#     },
#   },
#   {
#     $addFields: {
#       start_ts: {
#         $subtract: [
#           "$cst_ts",
#           {
#             $mod: ["$cst_ts", 3600 * 24],
#           },
#         ],
#       },
#     },
#   },
#   {
#     $addFields: {
#       start_date: {
#         $toDate: {
#           $multiply: ["$start_ts", 1000],
#         },
#       },
#     },
#   },
#   {
#     $group: {
#       _id: "$start_ts",
#       usage_record_count: {
#         $count: {},
#       },
#       records: {
#         $push: {},
#       },
#     },
#   },
#   {
#     $sort: {
#       _id: 1,
#     },
#   },
# ]
# 使用以上聚合管道，获得所有有记录的日期
result = colletion.aggregate(
    [
        {
            '$addFields': {
                'cst_ts': {
                    '$add': [
                        '$timestamp', 8 * 3600
                    ]
                }
            }
        }, {
            '$addFields': {
                'cst_date': {
                    '$toDate': {
                        '$multiply': [
                            '$cst_ts', 1000
                        ]
                    }
                }
            }
        }, {
            '$addFields': {
                'start_ts': {
                    '$subtract': [
                        '$cst_ts', {
                            '$mod': [
                                '$cst_ts', 3600 * 24
                            ]
                        }
                    ]
                }
            }
        }, {
            '$addFields': {
                'start_date': {
                    '$toDate': {
                        '$multiply': [
                            '$start_ts', 1000
                        ]
                    }
                }
            }
        }, {
            '$group': {
                '_id': '$start_ts',
                'usage_record_count': {
                    '$count': {}
                },
                'records': {
                    '$push': {}
                }
            }
        }, {
            '$sort': {
                '_id': 1
            }
        }
    ]
)

print("Aggregation done")

dates = []

for doc in result:
    # 保存时间戳
    dates.append(doc['_id'])
    
print("Found", len(dates), "dates")
    
# {
#     "begin": datetime对象,
#     "duration": 86400,
#     "usage_count": 当日qchatgpt-usage文档数量(timestamp为时间戳),
#     "active_host_count": 当日qchatgpt-usage文档中的remote_addr数量,
#     "new_host_count": 当日analysis_usage_remote_addrs文档中的remote_addr数量(created_at为时间戳),
# }
for d in dates:
    
    # 86400秒是一天
    begin = d - 8 * 3600
    duration = 86400
    
    datetime_object = datetime.datetime.fromtimestamp(begin)
    print("Calculating", datetime.datetime.fromtimestamp(d))
    
    # 当日qchatgpt-usage文档数量(timestamp为时间戳)
    usage_count = colletion.count_documents({
        'timestamp': {
            '$gte': begin,
            '$lt': begin + duration
        }
    })
    
    # 当日qchatgpt-usage文档中的remote_addr数量, 先以remote_addr为键分组，然后计数
    docs = colletion.aggregate(
        [
            {
                '$match': {
                    'timestamp': {
                        '$gte': begin,
                        '$lt': begin + duration
                    }
                }
            }, {
                '$group': {
                    '_id': '$remote_addr',
                    'count': {
                        '$count': {}
                    }
                }
            }
        ]
    )
    active_host_count = 0
    for doc in docs:
        active_host_count += 1
    
    # 当日analysis_usage_remote_addrs文档中的remote_addr数量(created_at为时间戳)
    new_host_count = db['analysis_usage_remote_addrs'].count_documents({
        'created_at': {
            '$gte': datetime.datetime.fromtimestamp(begin),
            '$lt': datetime.datetime.fromtimestamp(begin + duration)
        }
    })
    
    print({
        "begin": datetime.datetime.fromtimestamp(d),
        "duration": duration,
        "usage_count": usage_count,
        "active_host_count": active_host_count,
        "new_host_count": new_host_count,
    })
    
    # 存到analysis_daily中，先检查是否有相同begin和duration的文档，有则更新，无则插入
    if db['analysis_daily'].count_documents({
        'begin': datetime.datetime.fromtimestamp(d),
        'duration': duration
    }) > 0:
        db['analysis_daily'].update_one({
            'begin': datetime.datetime.fromtimestamp(d),
            'duration': duration
        }, {
            '$set': {
                'usage_count': usage_count,
                'active_host_count': active_host_count,
                'new_host_count': new_host_count,
            }
        })
    else:
        db['analysis_daily'].insert_one({
            'begin': datetime.datetime.fromtimestamp(d),
            'duration': duration,
            'usage_count': usage_count,
            'active_host_count': active_host_count,
            'new_host_count': new_host_count,
        })
    
    print("Done", datetime.datetime.fromtimestamp(d))