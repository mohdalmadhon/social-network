set -e

echo "Running database migrations..."
MIGRATE_OUT=$(migrate -path ./internal/migrations -database "sqlite3://./db/social_network.db?_foreign_keys=on" up 2>&1) || true
echo "$MIGRATE_OUT"


if echo "$MIGRATE_OUT" | grep -qi "error" && ! echo "$MIGRATE_OUT" | grep -qi "no change"; then
  echo "Migration failed, aborting startup."
  exit 1
fi

echo "Starting server..."
exec ./server
