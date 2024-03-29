#!/bin/bash
source local.env

export MIGRATION_DSN="host=pg-local port=5432 dbname=$POSTGRES_DBNAME user=$POSTGRES_USER password=$POSTGRES_PASSWORD sslmode=disable"

sleep 2 && goose -dir "${MIGRATION_DIR}" postgres "${MIGRATION_DSN}" up -v