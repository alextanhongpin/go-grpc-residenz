package main

import (
	"context"
	"log"
	"net/http"

	"github.com/golang/glog"

	gw "github.com/alextanhongpin/go-residenz/proto/listing"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
)

func run() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := gw.RegisterListingServiceHandlerFromEndpoint(ctx, mux, "localhost:9090", opts)
	if err != nil {
		return err
	}
	log.Println("listening to port *:8080")
	return http.ListenAndServe(":8080", mux)
}

func main() {
	defer glog.Flush()
	if err := run(); err != nil {
		glog.Fatal(err)
	}
}
