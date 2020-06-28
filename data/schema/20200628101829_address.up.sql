create table address
(
    balance    int8    not null default 0 check (balance >= 0),

    public_key bytea   not null unique,

    username   text    not null unique,

    frozen     boolean not null default false
);
