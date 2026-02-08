#!/bin/sh
set -e

  echo "🎭 Bentext - Starting up..."
echo "=========================="

if [ -z "$PG_URL" ]; then
    echo "⚠️  WARNING: PG_URL environment variable is not set"
    echo "   Server will start in demo mode (no database)"
else
    echo "🔄 Running database migrations..."
    ./migrate

    if [ $? -ne 0 ]; then
        echo "❌ Migration failed, exiting..."
        exit 1
    fi

    echo "✅ Migrations completed successfully"
fi

echo "🚀 Starting Bentext server..."

exec ./main

