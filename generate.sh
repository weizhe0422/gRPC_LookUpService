#!bin/bash

protoc gRPC_Lookupservice/pb/lookup.proto --go_out=plugins=grpc:.