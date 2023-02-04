#!/bin/sh

TARGET=../../fmx-deployment/binaries

echo "copying ..."
cp -v fap $TARGET
cp -v ../conf/fap.toml $TARGET
echo "done"
