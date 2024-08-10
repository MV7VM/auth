create table roles(
    id serial primary key,
    role text
);

create table users(
    id bigserial primary key,
    login text,
    password_hash blob, 
    token text,
    mail text,
    roleID int references roles(id)
);