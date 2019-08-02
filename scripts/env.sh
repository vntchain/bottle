#!/bin/sh

set -e

if [ ! -f "scripts/env.sh" ]; then
    echo "$0 must be run from the root of the repository."
    exit 2
fi

# Create fake Go workspace if it doesn't exist yet.
workspace="$PWD/build/_workspace"
root="$PWD"
bottledir="$workspace/src/github.com/vntchain"
if [ ! -L "$bottledir/bottle" ]; then
    mkdir -p "$bottledir"
    cd "$bottledir"
    ln -s ../../../../../. bottle
    cd "$root"
fi

# Set up the environment to use the workspace.
GOPATH="$workspace"
export GOPATH

GOBIN="$PWD/build/bin"
export GOBIN

cd "$bottledir/bottle"
PWD="$bottledir/bottle"
# Launch the arguments with the configured environment.
exec "$@"


