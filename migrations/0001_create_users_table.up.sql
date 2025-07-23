create table roles
(
    id         serial
        primary key,
    code       text not null
        unique,
    name       text not null,
    created_at timestamp default now()
);

alter table roles
    owner to postgres;

create table users
(
    id         serial
        primary key,
    name       varchar(500)            not null
        constraint users_pk
            unique,
    password   varchar(500)            not null,
    created_at timestamp default CURRENT_TIMESTAMP,
    role_id    integer   default 2
        references roles,
    blocked    boolean   default false not null,
    email      varchar(255)
);

alter table users
    owner to postgres;

insert into roles (id, code, name)
values (1, 'admin', 'Админ'),
       (2, 'user', 'Пользователь');