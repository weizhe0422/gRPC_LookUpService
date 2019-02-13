#!bin/bash

protoc Lookupservice/ServerStreaming/pb/lookup.proto --go_out=plugins=grpc:.