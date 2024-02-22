-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
   id bigint primary key not null,
   name text,
   phone text,
   address text
);

CREATE TABLE IF NOT EXISTS done_orders (
     id serial primary key,
     user_id bigint references users(id) not null,
     start date,
     done date
);

CREATE TABLE IF NOT EXISTS current_orders (
    id serial primary key,
    user_id bigint references users(id) not null,
    start date
);

CREATE TABLE IF NOT EXISTS products (
  id    serial primary key,
  current_order_id int references current_orders(id),
  done_order_id int references done_orders(id),
  name  text,
  size  text,
  color text,
  text  text,
  img text,
  amount  int
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS products;
DROP TABLE IF EXISTS current_orders;
DROP TABLE IF EXISTS done_orders;
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
