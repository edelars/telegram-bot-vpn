
create schema IF NOT EXISTS vpn;

use vpn;

create table IF NOT EXISTS users
(
    id         int auto_increment,
    created_at timestamp default current_timestamp not null,
    tg_login   varchar(256)                        null,
    referal_id varchar(256)                        null,
    invite_referal_id varchar(256)                 null,
    expired_at timestamp                           null,
    primary key (id)
);


create index users_tg_login_index
    on users (tg_login);


create index users_referal_id_index
    on users (referal_id);



