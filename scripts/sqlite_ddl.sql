pragma foreign_keys = ON;

create table client
(
    id      text primary key,
    name    text not null,
    api_key text not null,
    status  int not null default 0
);