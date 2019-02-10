package main

import (
	"fmt"
	"github.com/weizhe0422/gRPC/LookupService/Unary/pb"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
)

func main() {
	var (
		conn   *grpc.ClientConn
		err    error
		client pb.LookUpServiceClient
		firstName, lastName string
	)
	if conn, err = grpc.Dial(":50051", grpc.WithInsecure()); err != nil {
		log.Fatalf("failed to coneect to gRPC server: %v", err)
	}

	defer conn.Close()

	client = pb.NewLookUpServiceClient(conn)

	firstName = "WeiZhe"
	lastName = "Chang"

	doUnary(client, firstName, lastName)
}

func doUnary(client pb.LookUpServiceClient, firstName, lastName string){
	var(
		lookUpReq *pb.LookUpRequest
		lookUpResp *pb.LookUpResponse
		err error
	)

	lookUpReq = &pb.LookUpRequest{
		LookingUp: &pb.LookingUp{
			FirstName:firstName,
			LastName:lastName,
		},
	}

	if lookUpResp, err = client.LookUp(context.Background(),lookUpReq); err!=nil{
		log.Fatalf("failed to send LookUp: %v", err)
	}

	fmt.Println("LookUp result: ", lookUpResp.Result)
}