#!/bin/bash

set -eu -o pipefail

EXE=~/dev/go/go_build_oldest

TEMPDIR=$(mktemp -d)

testdir_cleanup() {
    set +u
    if [[ "${TESTDEBUG}x" = "x" ]]; then
	rm -rf "$TEMPDIR"
    fi
}

# clean up on exit
#trap 'rm -rf "$TEMPDIR"' EXIT
trap 'testdir_cleanup' EXIT

# set up env in TEMPDIR
echo "Running tests in $TEMPDIR"

cp $EXE "$TEMPDIR/oldest"
cp ./*.test "$TEMPDIR/"

tar -xf testdir.tar -C "$TEMPDIR"

# run tests
cd "$TEMPDIR"
shelltest -c "$PWD"
