-- +goose Up
create table if not exists users (
    id bigint generated always as identity primary key,
    user_uuid uuid not null unique default gen_random_uuid(),
    info jsonb not null,
    created_at timestamp not null default current_timestamp,
    updated_at timestamp,
    password_hash text not null
);

create index idx_users_user_uuid on users(user_uuid);
create index idx_users_created_at on users(created_at);

-- +goose Down
drop table if exists users;
