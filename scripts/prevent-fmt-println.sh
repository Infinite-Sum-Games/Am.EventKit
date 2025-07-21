#!/bin/sh

if grep -nE 'fmt\.Print(ln|f)' "$@" ; then
  echo "Error: Found 'fmt.Println' or 'fmt.Printf' in staged Go files."
  echo "Please use structured logging (zerolog) instead."
  exit 1
fi
