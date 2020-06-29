create table operation
(
    consensus      int8 primary key not null,
    signature      bytea unique     not null,

    operation      text             not null,

    from_address   bytea,
    to_address     bytea,

    amount         int8             not null default 0,

    status         text             not null,
    status_message text             not null
);
