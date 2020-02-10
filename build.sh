#!/usr/bin/env bash

env GOOS=linux GOARCH=arm GOARM=7 go build -o build/mrsans github.com/skyzh/MrSans/mrsans
