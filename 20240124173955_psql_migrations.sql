-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    id bigint primary key not null,
    name text,
    phone text,
    address text
);

CREATE TABLE done_orders (
    id serial primary key,
    user_id bigint references users(id) not null,
    start date,
    done date
);

CREATE TABLE current_orders (
    id serial primary key,
    user_id bigint references users(id) not null,
    start date
);

CREATE TABLE products (
    id    serial primary key,
    current_order_id int references current_orders(id),
    done_order_id int references done_orders(id),
    name  text,
    size  text,
    color text,
    text  text,
    img   text,
    amount  int
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE products;
DROP TABLE done_orders;
DROP TABLE current_orders;
DROP TABLE users;
-- +goose StatementEnd
