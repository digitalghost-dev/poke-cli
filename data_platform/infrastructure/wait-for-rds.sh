#!/bin/bash

MAX_TRIES=20
COUNT=0

RDS_HOST="<rds-instance-id>.<region>.rds.amazonaws.com"
RDS_PORT=5432

echo "Checking if RDS is available..."

while [ $COUNT -lt $MAX_TRIES ]; do
    if nc -z -w5 $RDS_HOST $RDS_PORT 2>/dev/null; then
        echo "RDS is available!"
        exit 0
    fi
    COUNT=$((COUNT + 1))
    echo "Attempt $COUNT/$MAX_TRIES - RDS not ready yet..."
    sleep 10
done

echo "RDS did not become available in time"
exit 1