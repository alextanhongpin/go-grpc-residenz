package main

import (
	"flag"
	"log"
	"net"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/alextanhongpin/go-residenz/app/database"
	pb "github.com/alextanhongpin/go-residenz/proto/listing"
)

type listingserver struct {
	DB *database.Database
}

// The database collection name
const collection string = "listing"

// ListingWithID contains the bson.ObjectId that cannot be parsed by protobuf
type ListingWithID struct {
	ID         bson.ObjectId    `bson:"_id,omitempty" json:"-"` // Map the _id field to golang struct bson.ObjectId
	pb.Listing `bson:",inline"` // If we did not use inline, it will be saved as a nested object in mongodb
}

func (s *listingserver) GetListings(ctx context.Context, msg *pb.GetListingsRequest) (*pb.GetListingsResponse, error) {
	session := s.DB.CopySession()
	defer session.Close()
	c := s.DB.Collection(session, collection)

	var listings []ListingWithID
	if err := c.Find(bson.M{}).All(&listings); err != nil {
		return nil, err
	}

	var responses []*pb.Listing
	for _, v := range listings {
		responses = append(responses, &pb.Listing{
			Id:          v.ID.Hex(),
			CreatedAt:   v.CreatedAt,
			ModifiedAt:  v.ModifiedAt,
			Title:       v.Title,
			Cost:        v.Cost,
			Description: v.Description,
			CoverPhoto:  v.CoverPhoto,
			Address:     v.Address,
			IsPublished: v.IsPublished,
			IsAvailable: v.IsAvailable,
		})
	}

	return &pb.GetListingsResponse{Data: responses}, nil
}

func (s *listingserver) GetListing(ctx context.Context, msg *pb.GetListingRequest) (*pb.GetListingResponse, error) {
	session := s.DB.CopySession()
	defer session.Close()
	c := s.DB.Collection(session, collection)

	var listing ListingWithID
	if err := c.FindId(bson.ObjectIdHex(msg.Id)).One(&listing); err != nil {
		return nil, err
	}

	data := &pb.Listing{
		Id:          listing.ID.Hex(), // Manual conversion of objectID to string
		CreatedAt:   listing.CreatedAt,
		ModifiedAt:  listing.ModifiedAt,
		Title:       listing.Title,
		Cost:        listing.Cost,
		Description: listing.Description,
		CoverPhoto:  listing.CoverPhoto,
		Address:     listing.Address,
		IsPublished: listing.IsPublished,
		IsAvailable: listing.IsAvailable,
	}

	return &pb.GetListingResponse{
		Data: data,
	}, nil
}

func (s *listingserver) UpdateListing(ctx context.Context, msg *pb.UpdateListingRequest) (*pb.UpdateListingResponse, error) {
	session := s.DB.CopySession()
	defer session.Close()

	c := s.DB.Collection(session, collection)
	if err := c.UpdateId(bson.ObjectIdHex(msg.Id), &msg); err != nil {
		return nil, err
	}
	return &pb.UpdateListingResponse{
		Msg: "Successfully updated response",
	}, nil
}

func (s *listingserver) PostListing(ctx context.Context, msg *pb.PostListingRequest) (*pb.PostListingResponse, error) {
	session := s.DB.CopySession()
	defer session.Close()

	c := s.DB.Collection(session, collection)
	// TODO: Validate the payload
	msg.CreatedAt = time.Now().Unix()
	msg.ModifiedAt = time.Now().Unix()
	err := c.Insert(msg)
	if err != nil {
		if mgo.IsDup(err) {
			return nil, err
		}
		return nil, err
	}

	return &pb.PostListingResponse{
		Msg: "Successfully created listing",
	}, nil
}

func (s *listingserver) DeleteListing(ctx context.Context, msg *pb.DeleteListingRequest) (*pb.DeleteListingResponse, error) {
	session := s.DB.CopySession()
	defer session.Close()

	c := s.DB.Collection(session, collection)
	err := c.RemoveId(bson.ObjectIdHex(msg.Id))
	if err != nil {
		return nil, err
	}
	return &pb.DeleteListingResponse{
		Msg: "Successfully deleted listing",
	}, nil
}

func main() {
	var (
		mongoHost     = flag.String("MONGO_HOST", "127.0.0.1:27017", "The mongodb host address")
		mongoUsername = flag.String("MONGO_USERNAME", "", "The mongodb username")
		mongoPassword = flag.String("MONGO_PASSWORD", "", "The mongodb password")
		mongoDatabase = flag.String("MONGO_DATABASE", "go-grpc-residenz", "The mongodb database name")
		port          = flag.String("PORT", ":9090", "The port the server is listening to")
	)
	flag.Parse()

	db, err := database.New(database.Config{
		Addrs:    []string{*mongoHost},
		Username: *mongoUsername,
		Password: *mongoPassword,
		Database: *mongoDatabase,
	})
	defer db.Session.Close()
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}

	log.Println("connected to database")

	lis, err := net.Listen("tcp", *port)
	if err != nil {
		log.Fatalf("failed to listen %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterListingServiceServer(grpcServer, &listingserver{DB: db})
	log.Printf("listening to port *%s\n", *port)
	grpcServer.Serve(lis)
}
