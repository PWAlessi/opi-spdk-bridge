// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2022 Dell Inc, or its subsidiaries.

package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	pb "github.com/opiproject/opi-api/storage/v1alpha1/gen/go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	port = flag.Int("port", 50051, "The Server port")
)

type Server struct {
	pb.UnimplementedFrontendNvmeServiceServer
	pb.UnimplementedNVMfRemoteControllerServiceServer
	pb.UnimplementedFrontendVirtioBlkServiceServer
	pb.UnimplementedFrontendVirtioScsiServiceServer
	pb.UnimplementedNullDebugServiceServer
	pb.UnimplementedAioControllerServiceServer
	pb.UnimplementedMiddleendServiceServer
}

var opiSpdkServer Server

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	pb.RegisterFrontendNvmeServiceServer(s, &Server{})
	pb.RegisterNVMfRemoteControllerServiceServer(s, &Server{})
	pb.RegisterFrontendVirtioBlkServiceServer(s, &Server{})
	pb.RegisterFrontendVirtioScsiServiceServer(s, &Server{})
	pb.RegisterNullDebugServiceServer(s, &Server{})
	pb.RegisterAioControllerServiceServer(s, &Server{})
	pb.RegisterMiddleendServiceServer(s, &Server{})

	reflection.Register(s)

	log.Printf("Server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
