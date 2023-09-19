#!/bin/bash

# Fetch SSH connection parameters from environment variables
SSH_KEY="$SSH_KEY_ENV"               # SSH private key
SSH_USER="$SSH_USER_ENV"             # SSH user
SSH_HOST="$SSH_HOST_ENV"             # EC2 instance IP or hostname
REMOTE_DIR="$REMOTE_DIR_ENV"         # Remote Git repository directory

# Ensure that required environment variables are set
if [ -z "$SSH_KEY" ] || [ -z "$SSH_USER" ] || [ -z "$SSH_HOST" ] || [ -z "$REMOTE_DIR" ]; then
  echo "One or more required environment variables are not set."
  exit 1
fi

# SSH into the EC2 instance and update the code
ssh -i "$SSH_KEY" "$SSH_USER"@"$SSH_HOST" << EOF
  cd "$REMOTE_DIR"
  git pull

  # Build your Go application (replace with the actual build command)
  make build
  nohup ./bin/app > /dev/null 2>&1 &
EOF