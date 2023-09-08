import os
import sys

MYSQL_HOST = os.environ.get('MYSQL_HOST', 'localhost')
MYSQL_PORT = os.environ.get('MYSQL_PORT', '3306')
MYSQL_USER = os.environ.get('MYSQL_USER', 'root')
MYSQL_PASSWORD = os.environ.get('MYSQL_PASSWORD', 'password')
MYSQL_DB = os.environ.get('MYSQL_DB', 'test')

import pymysql

conn = pymysql.connect(host=MYSQL_HOST, port=int(MYSQL_PORT), user=MYSQL_USER, passwd=MYSQL_PASSWORD, db=MYSQL_DB)

print("Connected to MySQL")

cursor = conn.cursor()

cursor.execute("select count(*) from `usage`")

print("MySQL query done")

print("Total", cursor.fetchone()[0])