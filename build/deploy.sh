#!/usr/bin/env bash

cd ..

git stash
git pull

go build cmd/marketplace/main.go
./main &
