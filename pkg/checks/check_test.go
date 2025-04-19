package checks_test

import (
	"slices"
	"testing"

	_ "github.com/mattn/go-sqlite3"

	"sqlite_lint/pkg/checks"
)

func TestFailureCases(t *testing.T) {
	test_cases := []struct {
		sqlFile          string
		expectedFailures []string
	}{
		{"../../test_schemas/failure-has-foreign-key-no-index.sql", []string{"require_indexes_for_foreign_keys"}},
		{"../../test_schemas/failure-has-ints.sql", []string{"forbid_int_type"}},
		{"../../test_schemas/failure-has-nulls.sql", []string{"require_not_null"}},
		{"../../test_schemas/failure-no-strict.sql", []string{"require_strict"}},
		{"../../test_schemas/failure-total.sql", []string{
			"require_not_null",
			"require_explicit_primary_key",
			"forbid_int_type",
			"require_strict",
			"require_indexes_for_foreign_keys",
		}},
	}

	for _, test_case := range test_cases {
		db, err := checks.OpenSchema(test_case.sqlFile)
		if err != nil {
			t.Fatalf("failed to open database: %v", err)
		}

		for _, check := range checks.Checks {
			results, err := check.Execute(db)
			if err != nil {
				t.Errorf("failed to execute check '%s': %v", check.Name, err)
			}

			is_failure := len(results) > 0
			is_failure_expected := slices.Contains(test_case.expectedFailures, check.Name)

			if is_failure != is_failure_expected {
				if is_failure_expected {
					t.Errorf("Expected check '%s' to fail, but it passed: %s", check.Name, test_case.sqlFile)
				} else {
					t.Errorf("Expected check '%s' to pass, but it failed: %s", check.Name, test_case.sqlFile)
				}
			}
		}
	}
}

func TestSuccessCase(t *testing.T) {
	db, err := checks.OpenSchema("../../test_schemas/success.sql")
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}

	for _, check := range checks.Checks {
		results, err := check.Execute(db)
		if err != nil {
			t.Errorf("failed to execute check '%s': %v", check.Name, err)
		}
		if len(results) > 0 {
			t.Errorf("Should have passed, but didn't: %s", "test_schemas/success.sql")
		}
	}
}
