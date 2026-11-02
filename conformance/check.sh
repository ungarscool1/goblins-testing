#!/bin/sh
# Trusted: Goblins always runs this version from the base commit.
grep -q '"/health"' main.go || { echo "main.go must route /health"; exit 1; }
echo "conformance passed"
