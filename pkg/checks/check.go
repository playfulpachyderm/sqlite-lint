package checks

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

// Check represents a database check with a name, SQL statement, and an explanation.
type Check struct {
	Name        string
	Sql         string
	Explanation string
}

// CheckResult represents a row in the query result with error message, table name, and column name.
type CheckResult struct {
	ErrorMsg   string `db:"error_msg"`
	TableName  string `db:"table_name"`
	ColumnName string `db:"column_name"`
}

// Execute runs the SQL statement of the Check against the provided database and returns the
// resulting rows as a slice of CheckResult using sqlx.
func (c *Check) Execute(db *sqlx.DB) ([]CheckResult, error) {
	var results []CheckResult
	// return results, nil
	// println(c.Sql)
	if err := db.Select(&results, c.Sql); err != nil {
		return nil, fmt.Errorf("failed to execute check '%s': %w", c.Name, err)
	}
	return results, nil
}

var Checks = map[string]Check{
	"require_not_null": {
		Name: "require_not_null",
		Sql: `
		    select 'Column should should be "not null"' as error_msg,
		           table_name,
		           column_name
		      from columns
		     where columns."notnull" = 0
		       and fk_target_column is null
		       and is_primary_key = 0 -- primary keys are automatically not-null, but aren't listed as such in pragma_table_info
		`,
		Explanation: "All columns should be marked as `not null` unless they are foreign keys.  (Primary keys are\n" +
			"automatically not-null, and don't need to be specified.)",
	},
	// {
	// 	Name: "require_default_values",
	// 	Sql: `
	// 	    select 'Column should have a default value' as error_msg,
	// 	           table_name,
	// 	           column_name
	// 	      from columns
	// 	     where dflt_value is null
	// 	       and fk_target_column is null
	// 	       and is_primary_key = 0;
	// 	`,
	// 	Explanation: "All columns should have a default value specified, unless they are foreign keys or primary keys.",
	// },
	"require_strict": {
		Name: "require_strict",
		Sql: `
		    select 'Table should be marked "strict"' as error_msg,
		           name as table_name,
		           '' as column_name
		      from tables
		     where strict = 0;
		`,
		Explanation: "All tables should be marked as `strict` (must specify column types; types must be int,\n" +
			"integer, real, text, blob or any).  This disallows all 'date' and 'time' column types.\n" +
			"See more: https://www.sqlite.org/stricttables.html",
	},
	"forbid_int_type": {
		Name: "forbid_int_type",
		Sql: `
		    select 'Column should use "integer" type instead of "int"' as error_msg,
		           table_name,
		           column_name
		      from columns
		     where column_type like 'int';
		`,
		Explanation: "All columns should use `integer` type instead of `int`.",
	},
	"require_explicit_primary_key": {
		Name: "require_explicit_primary_key",
		Sql: `
		    select 'Table should declare an explicit primary key' as error_msg,
		           tables.name as table_name,
		           '' as column_name
		      from tables
		     where not exists (select 1 from pragma_table_info(tables.name) where pk != 0);
		`,
		Explanation: "All tables must have a primary key.  If it's rowid, it has to be named explicitly.",
	},
	"require_indexes_for_foreign_keys": {
		Name: "require_indexes_for_foreign_keys",
		Sql: `
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
			           and is_primary_key != 0 -- 'pk' is either 0, or the 1-based index of the column within the primary key
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
		`,
		Explanation: "Columns referenced by foreign keys must have indexes.",
	},
}
