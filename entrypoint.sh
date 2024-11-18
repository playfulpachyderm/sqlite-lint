#!/bin/sh -l

set -x
set -e

if [ -z "$1" ]; then
	echo "No SQL schema file given!  Exiting..."
	exit 1
fi

DB_PATH=/tmp/database.db
SCHEMA_PATH="$1"
pwd
echo $SCHEMA_PATH

# Create the database
sqlite3 $DB_PATH < $SCHEMA_PATH

output=$(sqlite3 -column -header $DB_PATH < /lints.sql)
if [ -n "$output" ]; then
	echo "Some checks failed."
	echo
	echo $output
	exit 2
fi
