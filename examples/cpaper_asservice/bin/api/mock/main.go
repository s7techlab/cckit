package main

import (
	"context"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/s7techlab/cckit/examples/cpaper_asservice"
	cpaperservice "github.com/s7techlab/cckit/examples/cpaper_asservice/service"
	"github.com/s7techlab/cckit/gateway"
	servicemock "github.com/s7techlab/cckit/gateway/service/mock"
	"github.com/s7techlab/cckit/testing"
	"google.golang.org/grpc"
)

const (
	chaincodeName = `cpaper`
	channelName   = `cpaper`

	grpcAddress = `:8080`
	restAddress = `:8081`
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Create mock for commercial paper chaincode invocation

	// Commercial paper chaincode instance
	cc, err := cpaper_asservice.NewCC()
	if err != nil {
		log.Fatalln(err)
	}

	// Mockstub for commercial paper
	cpaperMock := testing.NewMockStub(chaincodeName, cc)

	// Chaincode invocation service mock. For real network you can use example with hlf-sdk-go
	cpaperMockService := servicemock.New(testing.NewPeer().WithChannel(channelName, cpaperMock))

	// default identity for signing requests to peeer (mocked)
	apiIdentity, err := testing.IdentityFromFile(`MSP`, `../../../testdata/admin.pem`, ioutil.ReadFile)
	if err != nil {
		log.Fatalln(err)
	}
	// Generated gateway for access to chaincode from external application
	cpaperGateway := cpaperservice.NewCPaperGateway(
		cpaperMockService, // gateway use mocked chaincode access service
		channelName,
		chaincodeName,
		gateway.WithDefaultSigner(apiIdentity))

	grpcListener, err := net.Listen("tcp", grpcAddress)
	if err != nil {
		log.Fatalf("failed to listen grpc: %v", err)
	}

	// Create gRPC server
	s := grpc.NewServer()
	cpaperservice.RegisterCPaperServer(s, cpaperGateway)

	// Runs gRPC server in goroutine
	go func() {
		log.Printf(`listen gRPC at %s`, grpcAddress)
		if err := s.Serve(grpcListener); err != nil {
			log.Fatalf("failed to serve gRPC: %v", err)
		}
	}()

	// wait for gRPC service stared
	time.Sleep(3 * time.Second)

	// Register gRPC server endpoint
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err = cpaperservice.RegisterCPaperHandlerFromEndpoint(ctx, mux, grpcAddress, opts)
	if err != nil {
		log.Fatalf("failed to register handler from endpoint %v", err)
	}

	log.Printf(`listen REST at %s`, restAddress)

	// Start HTTP server (and proxy calls to gRPC server endpoint)
	if err = http.ListenAndServe(restAddress, mux); err != nil {
		log.Fatalf("failed to serve REST: %v", err)
	}
}
