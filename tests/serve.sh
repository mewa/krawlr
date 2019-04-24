#!/bin/sh

BASE=$(dirname $0)

cd "$BASE/$1"

echo Serving contents of $PWD
python3 -m http.server 8888
