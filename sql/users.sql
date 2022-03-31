
create schema IF NOT EXISTS vpn;

use vpn;

create table IF NOT EXISTS users
(
    id         int auto_increment,
    created_at timestamp default current_timestamp not null,
    tg_login   varchar(256)                        null,
    tg_id   bigint                      null,
    referal_id varchar(256)                        null,
    invite_referal_id varchar(256)                 null,
    expired_at timestamp                           null,
    used_test_period bool not null default false,
    password varchar(256) null,
    primary key (id),
    unique key (tg_login),
    unique key (tg_id)

);


create index users_tg_login_index
    on users (tg_login);


create index users_referal_id_index
    on users (referal_id);



