# chat

## Database

```sql
-- database
create database chat;

create table if not exists session (
  session_id varchar(36) not null primary key,
  user_id int not null foreign key references user (user_id)
  last_login int not null
);
```
