#!/bin/bash

# Get the current user's UID and GID
USER_ID=$(id -u)
GROUP_ID=$(id -g)

# Define the .env file path
ENV_FILE=".env"

# Create or overwrite the .env file with UID and GID
cat <<EOL > $ENV_FILE
# Auto-generated .env file
UID=$USER_ID
GID=$GROUP_ID
EOL

echo ".env file generated with UID=$USER_ID and GID=$GROUP_ID"
