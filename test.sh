#!/bin/sh

rm data/*
for file in test_schemas/failure-*; do
	echo "Testing '$file'"
	test -e output.txt && rm output.txt
	db_path="data/$(basename $file).db"
	sqlite3 $db_path < $file
	sqlite3 -column -header $db_path < lints.sql | tee output.txt
	if [ ! -s output.txt ]; then
		echo "Should have failed, but didn't: $file"
		exit 1
	fi
done

file="test_schemas/success.sql"
echo "Testing '$file'"
test -e output.txt && rm output.txt
db_path="data/$(basename $file).db"
sqlite3 $db_path < $file
sqlite3 -column -header $db_path < lints.sql | tee output.txt
if [ -s output.txt ]; then
	echo "Should have passed, but didn't: $file"
	exit 1
fi

echo "Tests passed!"
