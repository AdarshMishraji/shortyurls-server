#!/bin/bash

SERVER_ADDRESS="https://shortyurls-server.onrender.com/"

# Infinite loop to keep sending requests
while true
do
    # Send an HTTP GET request
    curl "$SERVER_ADDRESS"
    sleep 5  # Sleep for 1 second between requests
    # Optional: Add a sleep interval to avoid overwhelming the server
    # sleep 1  # Sleep for 1 second between requests
done
