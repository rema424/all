-- db
create database myproject;

-- user
create user devuser@localhost identified by 'Passw0rd!';

-- grant
grant all privileges on myproject.* to devuser@localhost;

-- table
create table if not exists person (
  id bigint auto_increment,
  name varchar(255),
  primary key (id)
);

create table if not exists favorite_food (
  id bigint auto_increment,
  user_id bigint,
  name varchar(255),
  primary key (id),
  unique (user_id, name),
  foreign key (user_id) references person (id) on update cascade on delete set null
);