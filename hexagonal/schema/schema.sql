create table User
(
    username     varchar(20) not null
        primary key,
    password     varchar(64) not null,
    display_name varchar(20) not null,
    created_at   bigint      not null,
    updated_at   bigint      not null
);