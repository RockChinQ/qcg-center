import os

import pymongo

source_uri = os.environ['MONGO_URI']
target_uri = os.environ['MONGO_URI_TARGET']

db_name = "qcg-center-records"

colletions = [
    # "qchatgpt-usage",
    # "installer-reports",
    "analysis_daily",
    "analysis_reports_remote_addrs",
    "analysis_usage_remote_addrs",
]

source_client = pymongo.MongoClient(source_uri)
target_client = pymongo.MongoClient(target_uri)

source_db = source_client[db_name]
target_db = target_client[db_name]

for collection in colletions:
    
    print("Processing collection", collection)
    
    source_collection = source_db[collection]
    target_collection = target_db[collection]
    
    # 每次5000条，全量迁移
    limit = 10000
    
    offset = 0
    
    while True:
        print("Processing offset", offset)
        
        cursor = source_collection.find({}).skip(offset).limit(limit)
        
        to_insert = []
        
        for doc in cursor:
            to_insert.append(doc)
            
        if len(to_insert) == 0:
            break
            
        target_collection.insert_many(to_insert)
        
        offset += limit
        
        print("Processed offset", offset)
        
    print("Processed collection", collection)
    
print("Done")