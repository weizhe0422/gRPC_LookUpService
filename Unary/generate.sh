#!bin/bash

protoc Lookupservice/Unary/pb/lookup.proto --go_out=plugins=grpc:.