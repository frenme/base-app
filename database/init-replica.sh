#!/bin/bash
set -e

PGDATA="/var/lib/postgresql/data"

if [ "$(ls -A "$PGDATA" 2>/dev/null)" ]; then
  echo "Data directory $PGDATA is not empty, skip pg_basebackup..."
else
  echo "Data directory $PGDATA is not empty, skip pg_basebackup..."
  pg_basebackup \
    --pgdata="$PGDATA" \
    -R \
    --host=postgres-master \
    --port=5432
fi

echo "waiting for primary to connect..."
sleep 5

echo "Done, starting replica..."
chmod 700 "$PGDATA"

exec postgres