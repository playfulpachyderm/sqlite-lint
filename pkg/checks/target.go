package checks

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

// OpenSchema opens a SQLite database in memory, executes the schema against it, and adds some views
func OpenSchema(filepath string) (*sqlx.DB, error) {
	// Open a SQLite database in memory
	db, err := sqlx.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, fmt.Errorf("failed to open in-memory database: %w", err)
	}

	// Read the SQL file
	sqlBytes, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read SQL file: %w", err)
	}

	// Execute the SQL statements
	db.MustExec(string(sqlBytes))

	// Execute the SQL statements for creating views
	db.MustExec(`
		create view tables as
			select l.*
			from sqlite_schema s
		left join pragma_table_list l on s.name = l.name
			where s.type = 'table';


		create view columns as
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
	`)

	return db, nil
}
