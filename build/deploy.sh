#!/usr/bin/env bash

git stash
git pull

go build cmd/marketplace/main.go
