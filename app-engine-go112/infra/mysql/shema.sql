-- db
create database myproject;

-- user
create user devuser@localhost identified by 'Passw0rd!';

-- grant
grant all privileges on myproject.* to devuser@localhost;

-- table
create table if not exists user (
  id int auto_increment,
  name varchar(255),
  primary key (id)
);

create table if not exists favorite_food (
  id int auto_increment,
  user_id int,
  name varchar(255),
  primary key (id),
  unique (user_id, name),
  foreign key (user_id) references user (id) on update cascade on delete set null
);