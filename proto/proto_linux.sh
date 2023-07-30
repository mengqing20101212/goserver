#!/bin/bash
../tool/protoc-23.4-linux-x86_64/bin/protoc --proto_path=. --go_out=../protobuf  *.proto
../tool/protoc-23.4-linux-x86_64/bin/protoc --proto_path=. --go_out=../protobuf  table/*.proto