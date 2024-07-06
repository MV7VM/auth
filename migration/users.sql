create table roles(
                      id serial primary key,
                      role text
);

create table users(
                      id bigserial primary key,
                      login text,
                      password text,
                      roleID int references roles(id)
);