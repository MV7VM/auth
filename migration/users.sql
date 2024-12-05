create table roles(
    id serial primary key,
    role text,
    grade int
);

create table users(
    id bigserial primary key,
    login text,
    role_id int references roles(id),
    password text
);