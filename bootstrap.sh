#!/bin/sh

if command -v gogitversion > /dev/null 2>&1; then
    gogitversion -p main
    exit 0
elif ls ./gogitversion > /dev/null 2>&1; then
    ./gogitversion -p main
elif ! ls ./main_version.go > /dev/null 2>&1; then
    echo 'package main\nconst version = "temporary"' > main_version.go
fi
