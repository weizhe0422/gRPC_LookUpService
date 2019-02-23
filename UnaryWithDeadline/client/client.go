package main

import (
	"fmt"
	"github.com/weizhe0422/gRPC_LookUpService/pb"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"time"
)

func main() {
	var (
		conn                *grpc.ClientConn
		err                 error
		client              pb.LookUpServiceClient
		firstName, lastName string
	)
	if conn, err = grpc.Dial(":50051", grpc.WithInsecure()); err != nil {
		log.Fatalf("failed to coneect to gRPC server: %v", err)
	}

	defer conn.Close()

	client = pb.NewLookUpServiceClient(conn)

	firstName = "WeiZhe"
	lastName = "Chang"

	doUnaryWithDeadline(client, firstName, lastName, 5*time.Second)
	doUnaryWithDeadline(client, firstName, lastName, 1*time.Second)
}

func doUnaryWithDeadline(client pb.LookUpServiceClient, firstName, lastName string, timeOut time.Duration) {
	var (
		lookUpReq  *pb.LookUpRequest
		lookUpResp *pb.LookUpResponse
		err        error
		ctx        context.Context
		cancel     context.CancelFunc
		ok         bool
		statusErr  *status.Status
	)

	lookUpReq = &pb.LookUpRequest{
		LookingUp: &pb.LookingUp{
			FirstName: firstName,
			LastName:  lastName,
		},
	}

	ctx, cancel = context.WithTimeout(context.Background(), timeOut)
	defer cancel()

	if lookUpResp, err = client.LookUpWithDeadline(ctx, lookUpReq); err != nil {
		statusErr, ok = status.FromError(err)
		if ok {
			if statusErr.Code() == codes.DeadlineExceeded {
				fmt.Println("Timeout was hit! Deadline was exceeded")
			} else {
				fmt.Printf("unexpected err: %v", statusErr)
			}
		} else {
			log.Fatalf("failed while calling LookUpWithDeadline RPC: %v", err)
		}
		return
	}

	fmt.Println("LookUp result: ", lookUpResp.Result)
}
