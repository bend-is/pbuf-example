#!/bin/bash
export PATH="$PATH:/protoc/bin"
BUFDIR=$(mktemp -d -t buf-XXXXXXXXXX)
buf build -o image.bin
protoc --descriptor_set_in=image.bin --go_out=$BUFDIR --go-grpc_out=$BUFDIR $(buf ls-files --input image.bin)
cp -r $BUFDIR/github.com/bend-is/pbuf-example/* ./
rm -rf $BUFDIR
rm -f image.bin
