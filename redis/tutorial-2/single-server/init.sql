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
  index (last_login)
);

create table room (
  id bigint auto_increment,
  name varchar(255),
  primary key (id)
);

create table room_user_rel (
  room_id bigint,
  user_id bigint,
  last_used_at timestamp default current_timestamp,
  primary key (room_id, user_id),
  foreign key (room_id) references room (id) on delete cascade on update cascade,
  foreign key (user_id) references user (id) on delete cascade on update cascade,
  index (user_id, last_used_at)
);

create table message (
  id bigint auto_increment,
  body varchar(255),
  sender_id bigint,
  room_id bigint,
  created_at timestamp default current_timestamp,
  primary key (id),
  foreign key (sender_id) references user (id)  on delete set null on update cascade,
  foreign key (room_id) references room (id)  on delete set null on update cascade,
  index (room_id, created_at)
);