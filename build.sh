#!/bin/bash

cd generate/main
go build .
 cd ../..
 generate/main/generate
 mv assets_vfsdata.go main
 cd main
 go build ./...
 mv ./main ../asset-hoster