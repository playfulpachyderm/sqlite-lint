package main

import (
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"

	"sqlite_lint/pkg/checks"
)

const (
	GREEN = "\033[0;32m"
	RED   = "\033[0;31m"
	RESET = "\033[0m"
)

func main() {
	// Check if a filepath argument is provided
	if len(os.Args) < 2 {
		fmt.Println("Please provide a filepath as the first argument.")
		os.Exit(1)
	}

	// Get the filepath from the first argument
	filepath := os.Args[1]

	fmt.Printf("-----------------\nLinting %s\n", filepath)

	// Open the SQLite database
	db, err := checks.OpenSchema(filepath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	is_failure := false
	// Execute each check against the database
	for _, check := range checks.Checks {
		results, err := check.Execute(db)
		if err != nil {
			panic(err)
		}

		// If there are results, print them as lint errors
		if len(results) > 0 {
			is_failure = true
			fmt.Printf(RED+"Check '%s' failed:\n"+RESET, check.Name)
			for _, result := range results {
				fmt.Printf(RED+"- %s: %s.%s\n"+RESET, result.ErrorMsg, result.TableName, result.ColumnName)
			}
			fmt.Printf(RED+"Explanation: %s\n\n"+RESET, check.Explanation)
		}
	}
	if is_failure {
		fmt.Println(RED + "Errors found" + RESET)
	} else {
		fmt.Println(GREEN + "Success" + RESET)
	}
}
