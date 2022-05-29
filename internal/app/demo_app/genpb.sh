#!/bin/bash

set -ex

dir="$(dirname "$0")"
path="$(cd "$dir";pwd -P)"
in="$path"/../../../proto/demo_app
out="$path"/../../../api/genproto/demo_app

sh "$path"/../../../scripts/genpb.sh "$in" "$out"

