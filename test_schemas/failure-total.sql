PRAGMA foreign_keys = on;


create table implicit_rowid (
    a integer
);

create table explicit_rowid_not_pk (
	rowid integer
);

create table explicit_rowid (
	rowid integer primary key
);

create table without_rowid (
	a integer primary key
) without rowid;

create table multi_column (
	a int,
	b integer,
	primary key(a, b)
);

create table foreign_key_missing_index (
	a references implicit_rowid(a)
);
