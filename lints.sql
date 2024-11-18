create temporary view tables as
    -- select * from sqlite_schema where type = 'table';
    select l.* from sqlite_schema s left join pragma_table_list l on s.name = l.name where s.type = 'table';


create temporary view columns as
    select tables.name as table_name,
           table_info.name as column_name,
           table_info.type as column_type,
           "notnull",
           dflt_value,
           pk as is_primary_key,
           fk."table" as fk_target_table,
           fk."to" as fk_target_column
      from tables
      join pragma_table_info(tables.name) as table_info
 left join pragma_foreign_key_list(tables.name) as fk on fk."from" = column_name;


-- ======================
-- Lint-check definitions
-- ======================

-- Require all fields to be not null, unless they are foreign keys
create temporary view lint_check__require_not_null as
    select 'Column should should be "not null"' as error_msg,
           tables.name as table_name,
           table_info.name as column_name
      from tables
      join pragma_table_info(tables.name) as table_info
     where not exists (
                select 1 from pragma_foreign_key_list(tables.name) as fk
                 where fk."from" = table_info.name
               )
       and (table_info."notnull" = 0 and table_info.pk = 0);


-- Require all non-foreign-key, non-primary-key fields to have a default value
create temporary view lint_check__require_default_values as
    select 'Column should have a default value' as error_msg,
           tables.name as table_name,
           table_info.name as column_name
      from tables
      join pragma_table_info(tables.name) as table_info
     where table_info.pk = 0
       and not exists (
                select 1 from pragma_foreign_key_list(tables.name) as fk
                 where fk."from" = table_info.name
               )
       and table_info.dflt_value is null;


-- All tables should be STRICT (must specify column types; types must be int, integer, real, text,
-- blob or any).  This disallows all 'date' and 'time' columns automatically.
-- See more: https://www.sqlite.org/stricttables.html
create temporary view lint_check__require_strict as
    select 'Table should be marked "strict"' as error_msg,
           name as table_name,
           '' as column_name
      from tables
     where strict = 0
       and name not in ('sqlite_schema', 'sqlite_temp_schema');


-- Forbid use 'int' column types
create temporary view lint_check__forbid_int_type as
    select 'Column should use "integer" type instead of "int"' as error_msg,
           tables.name as table_name,
           table_info.name as column_name
      from tables
      join pragma_table_info(tables.name) as table_info
     where table_info.type like 'int';


-- tables must have a primary key; if it's rowid, it has to be named explicitly
create temporary view lint_check__require_explicit_primary_key as
    select 'Table should declare an explicit primary key' as error_msg,
           tables.name as table_name,
           '' as column_name
      from tables
     where not exists (select name from pragma_table_info(tables.name) where pk != 0);


-- columns referenced by foreign keys must have indexes
create temporary view lint_check__require_indexes_for_foreign_keys as
with index_info as (
        select indexes.name as index_name,
               tables.name as table_name,
               columns.name as column_name,
               case when origin = 'c' then 'regular' when origin = 'u' then 'unique' when origin = 'pk' then 'primary key' else origin end as index_type
          from tables
          join pragma_index_list(tables.name) as indexes
          join pragma_index_info(indexes.name) as columns

         union

        select '[auto-generated rowid primary key index]' as index_name,
               tables.name as table_name,
               columns.name as column_name,
               'primary key' as index_type
          from tables
          join pragma_table_info(tables.name) as "columns"
         where columns.name = 'rowid'
           and columns.pk != 0 -- `pk` is either 0, or the 1-based index of the column within the primary key
    ), foreign_key_targets as (
        select tables.name as "fk_source_table",
               fk."from" as "fk_source_column",
               fk."table" as fk_target_table,
               fk."to" as fk_target_column
          from tables
          join pragma_foreign_key_list(tables.name) as fk
    )
select fk_source_table, fk_source_column, fk_target_table, fk_target_column, ifnull(index_info.table_name, 'NULL') as index_table, ifnull(index_info.column_name, 'NULL') as index_column
  from foreign_key_targets
  left join index_info on foreign_key_targets.fk_target_table = index_info.table_name and foreign_key_targets.fk_target_column = index_info.column_name
 where index_column = 'NULL';


-- ==============
-- Run the checks
-- ==============

-- select * from columns;

select * from (
    select * from lint_check__require_not_null
 -- union
 --    select * from lint_check__require_default_values
 union
    select * from lint_check__require_strict
 union
    select * from lint_check__forbid_int_type
 union
    select * from lint_check__require_explicit_primary_key
 union
    select 'Foreign keys should point to indexed columns' as error_msg,
           fk_source_table as table_name,
           fk_source_column as column_name
      from lint_check__require_indexes_for_foreign_keys
);

