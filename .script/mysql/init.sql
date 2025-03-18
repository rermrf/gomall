create database gomall_user;

use gomall_user;
create table gomall_user.user
(
    id        bigint auto_increment
        primary key,
    phone     varchar(191) null,
    email     varchar(191) null,
    password  longtext     not null,
    nick_name longtext     null,
    avatar    longtext     null,
    ctime     bigint       not null,
    utime     bigint       not null,
    constraint uni_users_email
        unique (email),
    constraint uni_users_phone
        unique (phone)
);




# 准备 canal 用户
CREATE USER 'canal'@'%' IDENTIFIED BY 'canal';
GRANT ALL PRIVILEGES ON *.* TO 'canal'@'%' WITH GRANT OPTION;