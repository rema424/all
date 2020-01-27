import sqlalchemy as sa

conn = sa.create_engine(
    "mysql+pymysql://tester:Passw0rd!@127.0.0.1/py_sqlalchemy_sandbox",
    encoding='utf8',
    echo=True)

res = conn.execute('''
CREATE TABLE IF NOT EXISTS zoo (
    critter VARCHAR(255) PRIMARY KEY,
    count INT,
    damages FLOAT
);
''')

print(res)
