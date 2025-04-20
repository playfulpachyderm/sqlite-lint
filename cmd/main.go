package main

import (
	"fmt"
	"os"
	"strings"

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
		fmt.Println(RED + "Please provide a filepath as the first argument." + RESET)
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
		// Checks can be disabled via Github config / environment variables
		if !is_check_enabled(check) {
			continue
		}
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
		os.Exit(1)
	}
	fmt.Println(GREEN + "Success" + RESET)
}

// github_actions_input_env_var converts an input name to the corresponding
// environment variable name used by GitHub Actions.
func github_actions_input_env_var(name string) string {
	// GitHub normalizes both hyphens and underscores to underscores, then uppercases the name
	normalized := strings.NewReplacer("-", "_", " ", "_").Replace(name)
	return "INPUT_" + strings.ToUpper(normalized)
}

// Setting the environment variable INPUT_REQUIRE_NOT_NULL="false" disables the "require_not_null" check
func is_check_enabled(c checks.Check) bool {
	val, is_set := os.LookupEnv(github_actions_input_env_var(c.Name))
	if !is_set {
		// Enable all checks by default
		return true
	}
	// Anything except "false" is true
	return val != "false"
}
