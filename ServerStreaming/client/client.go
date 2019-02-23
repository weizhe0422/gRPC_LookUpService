package main

import (
	"fmt"
	"github.com/weizhe0422/gRPC_LookUpService/pb"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"io"
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

	doServerStreaming(client, firstName, lastName)
}

func doServerStreaming(client pb.LookUpServiceClient, firstName, lastName string){
	var(
		lookUpReq *pb.LookUpRequest
		lookUpResp *pb.LookUpResponse
		lookUpStreaming pb.LookUpService_LookUpServerStreamingClient
		err error
	)

	lookUpReq = &pb.LookUpRequest{
		LookingUp: &pb.LookingUp{
			FirstName:firstName,
			LastName:lastName,
		},
	}


	if lookUpStreaming, err = client.LookUpServerStreaming(context.Background(),lookUpReq); err!=nil{
		log.Fatalf("failed to send LookUp: %v", err)
	}

	for{
		if lookUpResp, err = lookUpStreaming.Recv(); err!=nil{
			if err == io.EOF{
				log.Println("end of receive data")
				break
			}
			log.Println("err to receive data")
			return
		}

		fmt.Println("LookUp result: ", lookUpResp)
	}
}