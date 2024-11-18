create table stuff (
    rowid integer primary key,
    data text not null,
    amount integer not null
) strict;
create index index_stuff_amount on stuff (amount);

create table stuff2 (
    weird_pk integer primary key,
    label text not null unique,
    stuff_id integer references stuff(rowid),
    alternative_stuff_id integer references stuff(amount)
) strict;
