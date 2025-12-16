create table todos
(
    id          varchar(100) not null,
    title       varchar(100) not null,
    description text         null,
    user_id     varchar(100) not null,
    created_at  bigint       not null,
    updated_at  bigint       not null,
    deleted_at  bigint       null,
    primary key (id),
    foreign key (user_id) references users (id)
);
