create table users (
    id bigserial primary key,
    username varchar(50) unique not null,
    email varchar(50) unique not null,
    password varchar(50) not null,
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);
