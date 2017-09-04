package main

import (
	"log"
	"net"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	pb "github.com/alextanhongpin/go-residenz/proto/listing"
)

type listingserver struct{}

func (s *listingserver) GetListing(ctx context.Context, msg *pb.GetListingRequest) (*pb.GetListingResponse, error) {
	return &pb.GetListingResponse{
		Data: &pb.Listing{
			Id:          "1",
			CreatedAt:   1,
			UpdatedAt:   2,
			Name:        "john",
			Description: "something",
			Cost:        11.0,
		},
	}, nil
}

func (s *listingserver) GetListings(ctx context.Context, msg *pb.GetListingsRequest) (*pb.GetListingsResponse, error) {
	return &pb.GetListingsResponse{
		Data: []*pb.Listing{},
	}, nil
}

func (s *listingserver) UpdateListing(ctx context.Context, msg *pb.UpdateListingRequest) (*pb.UpdateListingResponse, error) {
	return nil, nil
}

func (s *listingserver) PostListing(ctx context.Context, msg *pb.PostListingRequest) (*pb.PostListingResponse, error) {
	return &pb.PostListingResponse{
		Id: "1",
	}, nil
}

func (s *listingserver) DeleteListing(ctx context.Context, msg *pb.DeleteListingRequest) (*pb.DeleteListingResponse, error) {
	log.Println(msg)
	return &pb.DeleteListingResponse{
		Msg: "hello",
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":9090")
	if err != nil {
		log.Fatalf("failed to listen %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterListingServiceServer(grpcServer, &listingserver{})
	log.Println("listening to port *:9090")
	grpcServer.Serve(lis)
}
