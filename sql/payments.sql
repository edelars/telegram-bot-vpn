create schema IF NOT EXISTS vpn;

use vpn;

create table IF NOT EXISTS payments
(
    id         int auto_increment,
    user_id         int,
    created_at timestamp default current_timestamp not null,
    currency   varchar(3)                        not null default 'RUB',
    value int                       null,
    checked bool                    not null default false,
    primary key (id),
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);


create index payments_id_index
    on payments (user_id);

create index payments_checked_index
    on payments (checked);

