#!/bin/sh
echo "copying ..."
scp -P24 fap fimatrix@matrixworks.cn:~/programs/
scp -P24 ../conf/fap.toml fimatrix@matrixworks.cn:~/programs/
echo "done"
