#!/usr/bin/env bash

set -eo pipefail

# where am i?
me="$0"
me_home=$(dirname "$0")
me_home=$(cd "$me_home" && pwd)

# deps
DLV="dlv"

# parse arguments
args=$(getopt d $*)
set -- $args
for i; do
  case "$i"
  in
    -d)
      debug="true";
      shift;;
    --)
      shift; break;;
  esac
done

if [ ! -z "$debug" ]; then
  ETC="${me_home}/v1/internal/fixtures" "$DLV" test $* -- -test.v
else
  ETC="${me_home}/v1/internal/fixtures" go test -v $*
fi
