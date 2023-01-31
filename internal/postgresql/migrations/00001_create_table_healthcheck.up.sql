create table if not exists healthcheck
(
    id     serial primary key,
    time   timestamp default now(),
    url    varchar not null,
    result varchar not null
);