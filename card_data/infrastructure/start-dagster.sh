#!/bin/bash

# Fetch secrets from AWS Secrets Manager
SUPABASE_SECRETS=$(aws secretsmanager get-secret-value \
    --secret-id dagster/supabase-creds \
    --region us-west-2 \
    --query SecretString \
    --output text)

AWS_RDS_SECRETS_PW=$(aws secretsmanager get-secret-value \
    --secret-id '<secret-name>' \
    --region us-west-2 \
    --query SecretString \
    --output text)

AWS_RDS_SECRETS_HN=$(aws secretsmanager get-secret-value \
    --secret-id rds-hostname \
    --region us-west-2 \
    --query SecretString \
    --output text)

# Extract values
SUPABASE_PASSWORD=$(echo "$SUPABASE_SECRETS" | jq -r '.password')
export SUPABASE_PASSWORD

SUPABASE_USER=$(echo "$SUPABASE_SECRETS" | jq -r '.user')
export SUPABASE_USER

AWS_RDS_PASSWORD=$(echo "$AWS_RDS_SECRETS_PW" | jq -r '.password')
export AWS_RDS_PASSWORD

AWS_RDS_HOSTNAME=$(echo "$AWS_RDS_SECRETS_HN" | jq -r '.hostname')
export AWS_RDS_HOSTNAME

DAGSTER_HOME=/home/ubuntu/card_data/card_data/
export DAGSTER_HOME

# Activate the virtual environment
source /home/ubuntu/card_data/card_data/.venv/bin/activate

# Start Dagster
exec dg dev --host 0.0.0.0 --port 3000