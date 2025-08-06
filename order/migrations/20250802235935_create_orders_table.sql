-- +goose Up
create table orders (
    order_uuid uuid primary key default uuid_generate_v4(),
    user_uuid uuid not null,
    part_uuids text[] not null,
    total_price decimal(10, 2) not null,
    transaction_uuid uuid,
    payment_method integer,
    status integer not null default 0,
    created_at timestamp with time zone default now(),
    updated_at timestamp with time zone default now()
);

create index idx_orders_user_uuid on orders(user_uuid);
create index idx_orders_status on orders(status);
create index idx_orders_created_at on orders(created_at);

-- +goose Down
drop table if exists orders
