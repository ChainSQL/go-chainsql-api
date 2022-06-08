#!/bin/bash
#set -x

docker-compose -f ./build/build.yaml run --no-deps --rm go-sdk-build
