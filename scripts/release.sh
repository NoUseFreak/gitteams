#!/usr/bin/env bash

set -e 

# Get the parent directory of where this script is.
SOURCE="${BASH_SOURCE[0]}"
while [ -h "$SOURCE" ] ; do SOURCE="$(readlink "$SOURCE")"; done
DIR="$( cd -P "$( dirname "$SOURCE" )/.." && pwd )"

cd "$DIR"

curl -sL http://bit.ly/gh-get | BIN_DIR=/tmp PROJECT=NoUseFreak/letitgo bash

/tmp/letitgo $(git describe --tags --abbrev=0)

