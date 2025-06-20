#!/bin/bash

echo "🚀 Starting Faslah Application..."

IS_SERVER=$(grep -E '^IP_ADDRESS=' .env | awk -F= '{print $2}' | xargs)

if [[ -n "$IS_SERVER" ]]; then
  echo "📦 Stashing local changes..."
  git stash --include-untracked

  echo "📥 Pulling latest changes from GitHub..."
  git pull origin master
else
  echo "🐘 Starting PostgreSQL Docker service..."
  docker-compose up -d

  echo "🔍 Verifying if PostgreSQL is listening on port 5432..."
  sleep 3
  if ! sudo lsof -i :5432 | grep LISTEN >/dev/null; then
    echo "⚠️ PostgreSQL not listening on port 5432. Attempting docker-compose restart..."
    docker-compose restart
    sleep 3
    if ! sudo lsof -i :5432 | grep LISTEN >/dev/null; then
      echo "❌ PostgreSQL still not listening on port 5432. Exiting."
      exit 1
    else
      echo "✅ PostgreSQL is now listening on port 5432 after restart."
    fi
  else
    echo "✅ PostgreSQL is listening on port 5432."
  fi
fi

echo "📦 Tidying Go modules..."
go mod tidy

PID=$(pgrep -f ./faslah)
if [[ -n "$PID" ]]; then
  echo "🛑 Stopping existing server (PID: $PID)..."
  kill "$PID" > /dev/null 2>&1 || true
  sleep 1
fi

echo "🧹 Removing old binary..."
rm -f faslah

echo "🔨 Building the application..."
go build -o faslah cmd/*.go

echo "🏃 Starting the application in the background..."
setsid ./faslah > faslah.log 2>&1 &

if [[ -n "$IS_SERVER" ]] && git stash list | grep -q .; then
  echo "🎯 Trying to reapply stash..."
  if git stash pop --quiet; then
    echo "✅ Stash applied successfully!"
  else
    echo "❌ Conflict detected while applying stash. Dropping it..."
    git reset --hard
    git stash drop
  fi
else
  echo "🎯 No stash to apply."
fi

# ✅ Final message
echo "✅ faslah is running. View logs with: tail -f faslah.log"

# 🪵 Follow logs
tail -f faslah.log
