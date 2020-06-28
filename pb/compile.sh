#!/usr/bin/env bash
cd pb
protoc --proto_path=./ --go_out=./ --go_opt=paths=source_relative messages.proto
