#!/bin/bash
set -e

# Function to create a new database
function create_database() {
    local database=$1
    echo "Creating database: $database"
    psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
        CREATE DATABASE $database;
EOSQL
    echo "Database $database created successfully"
}

# Check if POSTGRES_MULTIPLE_DATABASES is set
if [ -n "$POSTGRES_MULTIPLE_DATABASES" ]; then
    echo "Multiple database creation requested"
    for db in $(echo $POSTGRES_MULTIPLE_DATABASES | tr ',' ' '); do
        if [ "$db" != "$POSTGRES_DB" ]; then
            create_database $db
        fi
    done
    echo "All databases created"
else
    echo "Single database mode, skipping additional database creation"
fi