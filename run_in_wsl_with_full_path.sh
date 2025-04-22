#!/bin/bash

# Set the current user in the config file
echo '{"CurrentUserName": "testuser"}' > ~/.gatorconfig.json

# Make sure the config file is readable
chmod 644 ~/.gatorconfig.json

# Print out the config file contents
echo "Config file contents:"
cat ~/.gatorconfig.json

# Run the Go application with the full path to Go
export PATH=$PATH:/usr/local/go/bin
go run . addfeed "Hacker News RSS" "https://hnrss.org/newest" 