#!/bin/bash

set -eu -o pipefail

EXE=~/dev/go/go_build_oldest_now

TEMPDIR=$(mktemp -d)

# clean up on exit
trap 'rm -rf "$TEMPDIR"' EXIT

# set up env in TEMPDIR
echo "Running tests in $TEMPDIR"

cp $EXE "$TEMPDIR/oldest"
cp ./*.test "$TEMPDIR/"

tar -xf testdir.tar -C "$TEMPDIR"

# run tests
cd "$TEMPDIR"
shelltest -c "$PWD"
