package main

import (
	"context"
	"fmt"
	"github.com/weizhe0422/gRPC_LookUpService/pb"
	"google.golang.org/grpc"
	"io"
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

	doBiDiStreaming(client, firstName, lastName)
}

func doBiDiStreaming(client pb.LookUpServiceClient, firstName, lastName string) {
	var (
		stream   pb.LookUpService_BiDiLookUpClient
		err      error
		waitChan chan struct{}
		reqSlice []*pb.LookUpRequest
		resp     *pb.LookUpResponse
	)

	if stream, err = client.BiDiLookUp(context.Background()); err != nil {
		log.Fatalf("failed to creating stream %v", err)
		return
	}

	waitChan = make(chan struct{})

	reqSlice = []*pb.LookUpRequest{
		&pb.LookUpRequest{
			LookingUp: &pb.LookingUp{
				FirstName: "Ray",
				LastName:  "Chang",
			},
		},
		&pb.LookUpRequest{
			LookingUp: &pb.LookingUp{
				FirstName: "WeiZhe",
				LastName:  "Chang",
			},
		},
		&pb.LookUpRequest{
			LookingUp: &pb.LookingUp{
				FirstName: "HuiLun",
				LastName:  "Su",
			},
		},
	}

	//send a bunch of message
	go func() {
		for _, req := range reqSlice {
			if err = stream.Send(req); err != nil {
				log.Fatalf("failed to send requese request %v", err)
			}
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()

	//receive a bunch of messsge
	go func() {
		for {
			if resp, err = stream.Recv(); err != nil {
				if err == io.EOF {
					break
				}
				log.Fatalf("failed while receiving: %v", err)
				break
			}
			fmt.Printf("Received: %v\n", resp.GetResult())
		}
		close(waitChan)
	}()

	//keep waiting till get signal from channel
	<-waitChan
}
