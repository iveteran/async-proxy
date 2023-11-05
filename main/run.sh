#!/bin/sh

if [ ! -e ./fap ]; then
  make
fi

./fap -c ../conf/fap.toml
