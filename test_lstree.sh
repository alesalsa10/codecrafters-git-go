#!/bin/bash
#
# Local test for `ls-tree --name-only`.
# Builds a sample tree with REAL git, then diffs real git's output
# against your program's output. Prints PASS/FAIL.
#
# Usage: ./test_lstree.sh

set -e

REPO="$(cd "$(dirname "$0")" && pwd)"
WORK="/tmp/lstree-test"

# 1. fresh scratch repo built with real git
rm -rf "$WORK"
mkdir -p "$WORK"
cd "$WORK"
git init -q
echo hello > file1
mkdir dir1 dir2
echo a > dir1/file_in_dir_1
echo b > dir1/file_in_dir_2
echo c > dir2/file_in_dir_3
git add -A

# 2. tree sha that real git wrote into .git/objects
TREE=$(git write-tree)
echo "tree sha: $TREE"
echo

# 3. compare (run your program from inside $WORK so its .git path resolves)
if diff <(git ls-tree --name-only "$TREE") <("$REPO/your_program.sh" ls-tree --name-only "$TREE"); then
  echo
  echo "PASS"
else
  echo
  echo "FAIL (left = expected / real git, right = yours)"
  exit 1
fi
