create temporary view tables as
    select l.*
      from sqlite_schema s
 left join pragma_table_list l on s.name = l.name
     where s.type = 'table';


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
           table_name,
           column_name
      from columns
     where columns."notnull" = 0
       and fk_target_column is null
       and is_primary_key = 0; -- primary keys are automatically not-null, but aren't listed as such in pragma_table_info


-- Require all non-foreign-key, non-primary-key fields to have a default value
create temporary view lint_check__require_default_values as
    select 'Column should have a default value' as error_msg,
           table_name,
           column_name
      from columns
     where table_info.dflt_value is null
       and fk_target_column is null
       and is_primary_key = 0;


-- All tables should be STRICT (must specify column types; types must be int, integer, real, text,
-- blob or any).  This disallows all 'date' and 'time' columns automatically.
-- See more: https://www.sqlite.org/stricttables.html
create temporary view lint_check__require_strict as
    select 'Table should be marked "strict"' as error_msg,
           name as table_name,
           '' as column_name
      from tables
     where strict = 0;


-- Forbid use 'int' column types
create temporary view lint_check__forbid_int_type as
    select 'Column should use "integer" type instead of "int"' as error_msg,
           table_name,
           column_name
      from columns
     where column_type like 'int';


-- tables must have a primary key; if it's rowid, it has to be named explicitly
create temporary view lint_check__require_explicit_primary_key as
    select 'Table should declare an explicit primary key' as error_msg,
           tables.name as table_name,
           '' as column_name
      from tables
     where not exists (select 1 from pragma_table_info(tables.name) where pk != 0);


-- columns referenced by foreign keys must have indexes
create temporary view lint_check__require_indexes_for_foreign_keys as
with index_info as (
        select tables.name as table_name,
               columns.name as column_name
          from tables
          join pragma_index_list(tables.name) as indexes
          join pragma_index_info(indexes.name) as columns

         union

        select table_name,
               column_name
          from columns
         where column_name = 'rowid'
           and is_primary_key != 0 -- `pk` is either 0, or the 1-based index of the column within the primary key
    ), foreign_keys as (
        select * from columns where fk_target_column is not null
    )
select 'Foreign keys should point to indexed columns' as error_msg,
       foreign_keys.table_name as table_name,
       foreign_keys.column_name as column_name
  from foreign_keys
  left join index_info on foreign_keys.fk_target_table = index_info.table_name
                      and foreign_keys.fk_target_column = index_info.column_name
 where index_info.column_name is null;


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
    select * from lint_check__require_indexes_for_foreign_keys
);

