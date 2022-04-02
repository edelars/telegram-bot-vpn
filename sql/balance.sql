create schema IF NOT EXISTS vpn;

use vpn;

create table IF NOT EXISTS balance
(
    id         int auto_increment,
    user_id         int not null ,
    debt int default 0,
    credit int default 0,
    primary key (id),
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
    );

create index balance_id_index
    on balance (id);

create index balance_user_id_index
    on balance (user_id);
