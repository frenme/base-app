#!/bin/bash
set -e

echo "Running pg_basebackup..."
pg_basebackup --pgdata=/var/lib/postgresql/data \
              -R \
              --host=postgres-master \
              --port=5432

echo "waiting for primary to connect..."
sleep 5

echo "Done, starting replica..."
chmod 700 /var/lib/postgresql/data

exec postgres