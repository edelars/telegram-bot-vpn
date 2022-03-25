create schema IF NOT EXISTS vpn;

use vpn;

create table IF NOT EXISTS payments
(
    user_id         int,
    created_at timestamp default current_timestamp not null,
    currency   varchar(3)                        not null default 'RUB',
    value FLOAT                       null,
    payed_at timestamp                           not null,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
    );


create index users_tg_login_index
    on payments (user_id);