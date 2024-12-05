create table roles(
    id serial primary key,
    role text,
    grade int
);
insert into roles(role,grade) values ('CLIENT', 0);
insert into roles(role,grade) values ('ADMIN', 999);

create table users(
    id bigserial primary key,
    login text,
    role_id int references roles(id),
    password text
);