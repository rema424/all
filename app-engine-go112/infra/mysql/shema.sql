-- db
create database myproject;

-- user
create user devuser@localhost identified by 'Passw0rd!';

-- grant
grant all privileges on myproject.* to devuser@localhost;


create table if not exists user (
  id int auto_increment,
  name varchar(255),
  primary key (id)
);

create table if not exists favorite_food (
  id int auto_increment,
  user_id int,
  primary key (id),
  foreign key (user_id) references user (id) on update cascade on delete set null
)