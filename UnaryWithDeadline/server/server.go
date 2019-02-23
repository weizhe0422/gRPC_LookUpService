package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/weizhe0422/gRPC_LookUpService/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"net"
	"time"
)

type MyDataTable struct {
	ID   int32
	Name string
}

var (
	recordTable []MyDataTable
)

func InitData() {

	recordTable = make([]MyDataTable, 0)

	recordTable = append(recordTable, MyDataTable{
		ID:   13579,
		Name: "WeiZhe",
	})
	recordTable = append(recordTable, MyDataTable{
		ID:   24680,
		Name: "Ray",
	})
}

type serviceSvr struct{}

func (s *serviceSvr) LookUp(ctx context.Context, req *pb.LookUpRequest) (response *pb.LookUpResponse, err error) {
	var (
		firstName  string
		lookUpResp *pb.LookUpResponse
		content    MyDataTable
		result     int32
	)

	log.Printf("LookUp function invoked by %v", req)
	log.Println(recordTable)

	firstName = req.LookingUp.FirstName
	result = -1
	for _, content = range recordTable {
		log.Println(content.Name, "/", firstName)
		if content.Name == firstName {
			result = content.ID
			break
		}
	}

	if result == -1 {
		return nil, errors.New("failed to find ID, please check first name again")
	}

	lookUpResp = &pb.LookUpResponse{
		Result: result,
	}

	return lookUpResp, nil
}

func (s *serviceSvr) LookUpServerStreaming(req *pb.LookUpRequest, stream pb.LookUpService_LookUpServerStreamingServer) error {
	var (
		firstName  string
		lookUpResp *pb.LookUpResponse
		content    MyDataTable
		result     int32
		idx        int
	)

	log.Printf("LookUp function invoked by %v", req)
	log.Println(recordTable)

	firstName = req.LookingUp.FirstName
	result = -1
	for _, content = range recordTable {
		log.Println(content.Name, "/", firstName)
		if content.Name == firstName {
			result = content.ID
			break
		}
	}

	if result == -1 {
		return errors.New("failed to find ID, please check first name again")
	}

	for idx = 1; idx < 10; idx++ {
		lookUpResp = &pb.LookUpResponse{
			Result: result,
		}
		stream.Send(lookUpResp)
		time.Sleep(1 * time.Second)
	}

	return nil
}

func (s *serviceSvr) BiDiLookUp(stream pb.LookUpService_BiDiLookUpServer) (err error) {
	var (
		firstName  string
		lookUpResq *pb.LookUpRequest
		sendErr    error
	)

	log.Printf("BiDiLookUp function invoked")
	log.Println(recordTable)

	for {
		lookUpResq, err = stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("failed to receive the message from client: %v", err)
		}
		firstName = lookUpResq.LookingUp.FirstName
		lookUpFromDataTable(firstName)

		sendErr = stream.Send(&pb.LookUpResponse{
			Result: lookUpFromDataTable(firstName),
		})

		if sendErr != nil {
			log.Fatalf("failed while sending data to client %v", sendErr)
			return sendErr
		}
	}

	return nil
}

func (s *serviceSvr) LookUpWithDeadline(ctx context.Context, req *pb.LookUpRequest) (*pb.LookUpResponse, error){
	var (
		firstName  string
		lookUpResp *pb.LookUpResponse
		content    MyDataTable
		result     int32
	)

	log.Printf("LookUpWithDeadline function invoked by %v", req)
	log.Println(recordTable)

	for i:= 0; i < 3 ; i++ {
		if ctx.Err() == context.Canceled {
			fmt.Printf("Client canceld the request")
			return nil, status.Error(codes.Canceled, " client canceled the request")
		}
		time.Sleep(1 * time.Second)
	}

	firstName = req.LookingUp.FirstName
	result = -1
	for _, content = range recordTable {
		log.Println(content.Name, "/", firstName)
		if content.Name == firstName {
			result = content.ID
			break
		}
	}

	if result == -1 {
		return nil, errors.New("failed to find ID, please check first name again")
	}

	lookUpResp = &pb.LookUpResponse{
		Result: result,
	}

	return lookUpResp, nil
}

func lookUpFromDataTable(firstName string) int32 {
	var (
		content MyDataTable
	)

	for _, content = range recordTable {
		log.Println(content.Name, "/", firstName)
		if content.Name == firstName {
			return content.ID
		}
	}

	return -1
}

func InitGRPCSvr() (err error) {
	var (
		listener net.Listener
		grpcSvr  *grpc.Server
	)

	if listener, err = net.Listen("tcp", ":50051"); err != nil {
		log.Fatalf("failed to create a listener: %v", err)
		return
	}

	grpcSvr = grpc.NewServer()
	pb.RegisterLookUpServiceServer(grpcSvr, &serviceSvr{})

	if err = grpcSvr.Serve(listener); err != nil {
		log.Fatalf("failed to serve gRPC: %v", err)
		return
	}

	log.Println("initial gRPC server success...")

	return nil
}

func main() {
	var (
		err error
	)

	// Initial faked data for look up
	InitData()

	//Initial gRPC server

	if err = InitGRPCSvr(); err != nil {
		panic("failed to initialize gRPC server")
	}

}
