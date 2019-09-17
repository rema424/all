-   <https://www.gonum.org/post/introtogonum/>
-   <https://godoc.org/gonum.org/v1/gonum/stat#MeanStdDev>
-   <https://blog.yudppp.com/posts/csv_fast_upload/>
-   <https://qiita.com/nanakrk965/items/7bb77d3ab0dda23d56d4>
-   <https://github.com/GoogleCloudPlatform/golang-samples/blob/master/appengine/go11x/cloudsql/cloudsql.go>
-   <https://github.com/jmoiron/sqlx>

```sh
mkdir market-rate

cd $_

touch README.md

touch main.go

vi main.go

go mod init $(basename $(pwd))

mysql.server start

mysql -u root

mysql -u root -p

create database market;

create user rema@localhost identified by 'xxxxxxxx';

grant all privileges on market.* to rema@localhost;

\q

mysql -u rema -p

use market;

create table prices (
    word varchar(30),
    price float,
    title varchar(150),
    url varchar(255),
    site varchar(30),
    PRIMARY KEY(word, url),
    KEY(word, price)
);

create table rates (
  id int auto_increment,
  parts_name varchar(30),
  memo varchar(100),
  narrow_rate float,
  sample_size int,
  skew float,
  stddev float,
  min float,
  lower  float,
  median float,
  mean float,
  upper float,
  max float,
  ratio float,
  PRIMARY KEY(id),
  KEY(parts_name, skew)
);

desc prices;

desc rates;

\q

echo '.envrc' >> .gitignore

echo 'export DSN="rema:your-password@tcp(127.0.0.1:3306)/market?parseTime=true&sql_mode=%27TRADITIONAL,NO_AUTO_VALUE_ON_ZERO,ONLY_FULL_GROUP_BY%27"' >> .envrc

direnv allow

echo $DSN

vi main.go

mv -f ~/Downloads/items_all_XX.csv ./items.csv
```
