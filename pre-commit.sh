#!/bin/sh

# pre-commit.sh, src: https://codeinthehole.com/tips/tips-for-using-a-git-pre-commit-hook/

git stash -q --keep-index
./run_tests.sh
RESULT=$?
git stash pop -q
[ $RESULT -ne 0 ] && exit 1
exit 0

