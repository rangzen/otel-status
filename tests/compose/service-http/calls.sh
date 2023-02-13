#!/bin/bash

echo "Testing HTTP calls:"

echo " --- Testing normal HTTP ---"
curl --write-out '| CODE:%{http_code}' localhost:8080
echo && echo

echo " --- Testing JSON with 1 second delay ---"
curl --write-out '| CODE:%{http_code}' localhost:8081
echo && echo

echo " --- Testing HTTP with problem ---"
curl --write-out '| CODE:%{http_code}' localhost:8082
echo