-- db
create database mychat;

-- user
create user workuser@localhost identified by 'Passw0rd!';

-- grant
grant all privileges on mychat.* to workuser@localhost;

-- table

create table user (
  id bigint auto_increment,
  name varchar(255),
  email varchar(255) character set latin1 collate latin1_bin,
  password_hash varchar(255) character set latin1 collate latin1_bin,
  created_at timestamp default current_timestamp,
  updated_at timestamp default current_timestamp on update current_timestamp,
  primary key (id),
  unique (email)
);

create table session (
  session_id varchar(255) character set latin1 collate latin1_bin,
  user_id bigint,
  last_logged_in_at timestamp default current_timestamp,
  primary key (session_id),
  foreign key (user_id) references user (id) on delete set null on update cascade,
  index (last_logged_in_at)
);

create table room (
  id bigint auto_increment,
  name varchar(255),
  primary key (id)
);

create table message (
  id bigint auto_increment,
  room_id bigint,
  user_id bigint,
  body varchar(255),
  created_at timestamp default current_timestamp,
  primary key (id),
  foreign key (user_id) references user (id)  on delete set null on update cascade,
  foreign key (room_id) references room (id)  on delete set null on update cascade,
  index (room_id, id)
);

-- seed

insert into room (name) values ('Room_1'), ('Room_2'), ('Room_3'), ('Room_4'), ('Room_5');
