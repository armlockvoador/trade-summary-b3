#!/bin/bash

if [ -z "$BASH" ]; then echo "Please run this script with bash"; exit 1; fi

SCRIPT_PATH=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
ROOT_PATH="$(cd $SCRIPT_PATH && cd ../../ && pwd)"
MOCKS_PATH="$SCRIPT_PATH/mock"

mkdir -p "$MOCKS_PATH"

echo "clearing existing mocks..."
find "$MOCKS_PATH" -mindepth 1 -maxdepth 1 -type d -exec rm -rv {} \;

echo "generating project mocks..."
cd $ROOT_PATH && mockery && cd -

echo "done."