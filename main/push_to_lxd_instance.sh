#!/bin/bash

TARGET=debian-12/home/fimatrix/programs/

echo "pushing ..."
lxc file push fap $TARGET
lxc file push ../conf/fap.toml $TARGET
echo "done"
