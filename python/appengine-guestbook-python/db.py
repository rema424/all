import sqlalchemy as sa

conn = sa.create_engine(
    "mysql+pymysql://tester:Passw0rd!@127.0.0.1/py_sqlalchemy_sandbox",
    encoding='utf8',
    echo=True)

# res = conn.execute('''
# CREATE TABLE IF NOT EXISTS zoo (
#     critter VARCHAR(255) PRIMARY KEY,
#     count INT,
#     damages FLOAT
# );
# ''')
# print(res)

# ins = 'INSERT IGNORE INTO zoo (critter, count, damages) VALUES (%s, %s, %s);'
# conn.execute(ins, 'duck', 10, 0.0)
# conn.execute(ins, 'bear', 2, 1000.0)
# conn.execute(ins, 'weasel', 1, 2000.0)
# rows = conn.execute('SELECT * FROM zoo')
# for row in rows:
#     print(row)

meta = sa.MetaData()
zoo = sa.Table('zoo', meta,
               sa.Column('critter', sa.String(255), primary_key=True),
               sa.Column('count', sa.Integer),
               sa.Column('damages', sa.Float),
               )
meta.create_all(conn)
conn.execute(zoo.insert(('bear', 2, 1000.0)))
conn.execute(zoo.insert(('weasel', 1, 2000.0)))
conn.execute(zoo.insert(('duck', 10, 0)))
result = conn.execute(zoo.select())
rows = result.fetchall()
print(rows)
