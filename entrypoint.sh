#!/bin/sh

echo "Running configFileUpdater.go..."
go run scripts/configFileUpdater.go
if [ $? -ne 0 ]; then
  echo "Failed to generate config. Exiting."
  exit 1
fi

echo "Starting load balancer..."
exec ./loadbalancer --config /app/config.json
