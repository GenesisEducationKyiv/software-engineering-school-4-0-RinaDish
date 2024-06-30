CREATE USER postgres SUPERUSER;

create table emails
(
    id             serial
        constraint users_pk
            primary key,
    email          varchar                 not null
        constraint emails_uniq
            unique
);