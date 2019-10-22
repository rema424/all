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

insert into user (name, email, password_hash) values
('alice', 'alice@example.com', ''),
('bob', 'bob@example.com', ''),
('carol', 'carol@example.com', ''),
('dave', 'dave@example.com', ''),
('ellen', 'ellen@example.com', ''),
('frank', 'frank@example.com', '');

select @alice := id from user where email = 'alice@example.com';
select @bob := id from user where email = 'bob@example.com';
select @carol := id from user where email = 'carol@example.com';
select @dave := id from user where email = 'dave@example.com';
select @ellen := id from user where email = 'ellen@example.com';
select @frank := id from user where email = 'frank@example.com';

select @room := id from room where name = 'Room_1';

insert into message (room_id, user_id, body) values
(@room, @alice, '能をつかんとする人、'),
(@room, @bob, '「よくせざらんほどは、なまじひに人に知られじ。うちうちよく習ひ得て、さし出でたらんこそ、いと心にくからめ」'),
(@room, @carol, 'と常に言ふめれど、'),
(@room, @dave, 'かく言ふ人、一芸も習ひ得ることなし。'),
(@room, @ellen, '未だ堅固かたほなるより、上手の中に交りて、毀り笑はるゝにも恥ぢず、つれなく過ぎて嗜む人、'),
(@room, @frank, '天性、その骨なけれども、'),
(@room, @alice, '道になづまず、濫りにせずして、年を送れば、'),
(@room, @bob, '堪能の嗜まざるよりは、終に上手の位に至り、'),
(@room, @carol, '徳たけ、人に許されて、双なき名を得る事なり。'),
(@room, @dave, '天下のものの上手といへども、'),
(@room, @ellen, '始めは、不堪の聞えもあり、無下の瑕瑾もありき。'),
(@room, @frank, 'されども、その人、道の掟正しく、これを重くして、放埒せざれば、'),
(@room, @alice, '世の博士にて、万人の師となる事、'),
(@room, @bob, '諸道変るべからず。'),
(@room, @carol, '徒然草　第百五十段');