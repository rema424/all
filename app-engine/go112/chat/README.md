# chat

## Database

```sh
mysql -uroot
```

```sql
-- database
create database chat;

grant all privileges on chat.* to workuser@localhost;

create table if not exists user (
  user_id int not null primary key auto_increment,
  name varchar(50) null,
  email varchar(255) null
);

create table if not exists session (
  session_id varchar(36) not null primary key,
  user_id int not null,
  last_login int not null,
  foreign key (user_id) references user (user_id)
);

create table if not exists credential (
  user_id int not null primary key,
  pasword_hash varchar(100) not null,
  foreign key (user_id) references user (user_id)
);
```
