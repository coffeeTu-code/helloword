#!/usr/bin/env bash

protoc -I . --go_out=paths=source_relative,plugins=grpc:./go/ *.proto
