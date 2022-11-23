#!/bin/sh
go mod init matrix.works/fmx-async-proxy

echo -en "\nreplace matrix.works/fmx-common => ../fmx-common" >> go.mod
