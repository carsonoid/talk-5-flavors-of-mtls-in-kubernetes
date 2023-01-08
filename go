#!/bin/bash

if grep -q '^/tmp/' <<<$PWD; then
  rsync -ruhv --exclude="*.go" $REPO_DIR/ $PWD > /dev/null
  cp -r $REPO_DIR/go.mod $PWD/go.mod > /dev/null
  rsync -ruhv $REPO_DIR/internal $PWD > /dev/null
fi

export GO111MODULE=auto

case $1 in
  build)
    $GO_ORIG_BIN $@ .
    ;;
  *)
    echo "Running $@"
    $GO_ORIG_BIN "$@"
    ;;
esac
