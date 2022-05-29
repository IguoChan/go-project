#!/bin/bash

set -ex

dir="$(dirname "$0")"
path="$(cd "$dir";pwd -P)"
in="$path"/../../../proto/example_app
out="$path"/../../../api/genproto/example_app

sh "$path"/../../../scripts/genpb.sh "$in" "$out"
