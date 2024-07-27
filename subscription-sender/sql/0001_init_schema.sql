CREATE USER postgres SUPERUSER;

create table messages
(
    id             serial
        -- constraint users_pk
            primary key,
    created_at timestamp default now() not null,
    event_id int,
    event_type varchar,
    emails varchar[] NOT NULL,
    sending_time timestamp not null,
    rate float not null,
    sent boolean default false not null
);

