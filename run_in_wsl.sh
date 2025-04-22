#!/bin/bash

# Set the current user in the config file
echo '{"CurrentUserName": "testuser"}' > ~/.gatorconfig.json

# Run the Go application
go run . addfeed "Hacker News RSS" "https://hnrss.org/newest" 