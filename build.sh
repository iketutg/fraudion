#!/bin/bash

echo "Building..."
go install
resbuild=$?

# Upload executable to test server?

exit $resbuild
