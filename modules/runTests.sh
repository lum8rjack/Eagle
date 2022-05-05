#!/bin/bash

echo "Starting the tests"
go test -v
sleep 1

echo ""
echo "Starting the fuzzing"
go test -fuzz FuzzCheckCIDR -fuzztime=60s -cpu=1
sleep 1
echo ""
go test -fuzz FuzzCheckPorts -fuzztime=60s -cpu=1
