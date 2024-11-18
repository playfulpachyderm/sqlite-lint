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

sqlite3 -column -header $DB_PATH < /lints.sql | tee output.txt
if [ -s output.txt ]; then
	echo "Some checks failed."
	exit 2
fi
