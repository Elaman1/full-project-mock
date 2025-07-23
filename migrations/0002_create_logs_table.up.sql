create table logs
(
    id         serial
        primary key,
    user_id    integer not null,
    action     text    not null,
    created_at timestamp default now()
);

alter table logs
    owner to postgres;

