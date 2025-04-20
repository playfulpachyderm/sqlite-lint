# SQLite Schema Lint

This GitHub Action lints SQLite schema files to enforce various constraints. It is designed to ensure that your SQLite schemas adhere to best practices.

## Inputs

- **`schema-file`**: (Required) The SQL schema file to lint.

## Usage

To use this action in your workflow, include a step in your job that uses `playfulpachyderm/sqlite-lint`. Below is an example of how to set it up in a GitHub Actions workflow file:

```yaml
jobs:
  lint:
    runs-on: ubuntu-latest

    steps:
      - name: Checkout code
        uses: actions/checkout@v4

      - name: Validate SQL schema
        uses: playfulpachyderm/sqlite-lint@v1.0.0
        with:
          schema-file: pkg/db/schema.sql
```

## Available Checks

All checks are enabled by default.  To turn one off, you can use:

```yaml
    with:
        - schema-file: [...]
        - require_not_null: false
```

This will disable the `require_not_null` check.

### `require_not_null`

Enforce that all columns should be marked as `not null` unless they are foreign keys.

**Explanation**: Nulls are a common source of unexpected bugs, because they're usually an invalid state but often get created by mistake (e.g., you forgot to set a value).  Explicitly disabling nulls prevents such mistakes.

If you need to track "unset" values, add an `is_X_valid` column to make it explicit.  Note that in many cases, `0` or `""` (empty string) are sufficient null values, and an `is_X_valid` flag might not even be required.

Foreign keys are exempt in this check because `null` is the only value the integrity checker will accept to represent "this row has no related item".

### `require_strict`

Enforce that all tables should be marked as `strict`.

**Explanation**: By default, SQLite is very loose with what it accepts, and basically doesn't enforce any type checking.   "Strict" tables disable this "looseness" and enforce that inserted values match the stated type of the column.

"Strict" tables also limit to a small number of column types: `int`, `integer`, `real`, `text`, `blob` or `any`.  To represent dates / times, use Unix epoch times in milliseconds, and convert to formatted dates (and timezones) only when displaying the value to a user.  This is the most portable and least bug-prone method to handle dates.

See more about "strict" tables in SQLite's documentation: <https://sqlite.org/stricttables.html>

### `forbid_int_type`

Enforce that all columns should use `integer` type instead of `int`.

**Explanation**: This is an extension of "strict" tables, which allow two redundant integer types, `integer` and `int`.  This check standardizes the types further, permitting only `integer`.

### `require_explicit_primary_key`

Enforce that all tables must have a primary key.  If the primary key is `rowid`, it must be declared explicitly.

**Explanation**: All tables need to have a primary key, and it should usually be `rowid`.  Declaring it explicitly improves the readability of the schema.

### `require_indexes_for_foreign_keys`

Enforce that columns referenced by foreign keys must have indexes.

**Explanation**: Foreign keys are usually used for `join`s.  Joining on un-indexed columns is very slow.  Ensuring that all foreign-key-referenced columns have indexes will greatly improve the performance of database operations.
