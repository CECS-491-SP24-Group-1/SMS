#!/bin/bash

# Determine the shell configuration file based on the current shell
if [[ "$SHELL" == *"/bash" ]]; then
	RC_FILE="$HOME/.bashrc"
elif [[ "$SHELL" == *"/zsh" ]]; then
	RC_FILE="$HOME/.zshrc"
else
	echo "Unsupported shell. This script supports Bash and Zsh."
	exit 1
fi

# Check if the shell configuration file exists
if [[ ! -f "$RC_FILE" ]]; then
	echo "Shell configuration file not found: $RC_FILE"
	exit 1
fi

# Define the line to be added
line_to_add="export PATH=/home/$(whoami)/go/bin:\$PATH"

# Check if the line already exists in the shell configuration file
if ! grep -qF "$line_to_add" "$RC_FILE"; then
	# Append the line to the shell configuration file
	echo "$line_to_add" >> "$RC_FILE"
	echo "Line added to $RC_FILE"
else
	echo "Line already exists in $RC_FILE"
fi

# Optionally, source the shell configuration file to apply the changes immediately
# Note: This step is optional and depends on whether you want the changes to take effect immediately
# source "$RC_FILE"
