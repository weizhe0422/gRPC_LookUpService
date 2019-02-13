package main

import (
	"context"
	"errors"
	"github.com/weizhe0422/gRPC_LookUpService/ServerStreaming/pb"
	"google.golang.org/grpc"
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
		idx int
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

	for idx = 1; idx < 10; idx++{
		lookUpResp = &pb.LookUpResponse{
			Result: result,
		}
		stream.Send(lookUpResp)
		time.Sleep(1 * time.Second)
	}

	return nil
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
