#!/bin/sh

migrate=$(migrate -database ${DB_DSN} -path internal/db/migrations up 2>&1)

if [[ $migrate == *"error"* || $migrate == *"timed out"* ]]; then
    echo "failed when running migrations: $migrate" 1>&2
    exit 1
fi

echo "successfully ran migrations: $migrate"
exit 0