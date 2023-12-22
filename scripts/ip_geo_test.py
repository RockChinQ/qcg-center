import os
import json
import time
import traceback

from motor import motor_asyncio
from pymongo import server_api
import requests


async def get_ip_info(ip):
    url = 'http://ip-api.com/json/{}?lang=zh-CN'.format(ip)
    response = requests.get(url)
    return response.json()

async def main():
    client = motor_asyncio.AsyncIOMotorClient(
        os.environ.get('MONGODB_URI'),
        serverSelectionTimeoutMS=1000,
        server_api=server_api.ServerApi('1'),
    )

    db = client['qcg-center-records']

    result = {}

    failed_ips = []

    limit = 500
    skip = 100

    failed = 0

    amount = 0

    try:
        # 从 analysis_usage_remote_addrs 集合中查询所有的记录
        async for record in db.analysis_usage_remote_addrs.find().limit(limit).skip(skip):
            try:
                # 获取 ip 地址
                ip = record['remote_addr']
                # 获取 ip 地址的信息
                ip_info = await get_ip_info(ip)

                if ip_info['status'] != 'success':
                    failed_ips.append(ip)
                    raise Exception(ip_info)
                
                result[ip] = ip_info
                amount += 1
                print(amount, ip, ip_info)
                failed = 0
            except Exception as e:
                failed += 1
                traceback.print_exc()
                if failed > 10:
                    break
            time.sleep(1)
    except Exception as e:
        traceback.print_exc()

    # 将结果写入文件
    with open('ip_info.json', 'w') as f:
        json.dump(result, f, indent=4, ensure_ascii=False)

if __name__ == '__main__':
    import asyncio
    asyncio.run(main())