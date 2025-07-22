#!/bin/sh

UNFORMATTED=$(gofmt -l "$@")
if [ -n "$UNFORMATTED" ]; then
  echo "Detected unformatted files. Running gofmt..."
  gofmt -w "$@"
  git add "$@"
  echo "gofmt applied and changes staged."
fi
