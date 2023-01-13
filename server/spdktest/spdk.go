package spdktest

import (
	"context"
	pb "github.com/opiproject/opi-api/storage/v1alpha1/gen/go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	"log"
	"net"
	"os"
)

func StartSpdkMockupServer() net.Listener {
	// start SPDK mockup server
	var rpcSock = "/var/tmp/spdk.sock"
	if err := os.RemoveAll(rpcSock); err != nil {
		log.Fatal(err)
	}
	ln, err := net.Listen("unix", rpcSock)
	if err != nil {
		log.Fatal("listen error:", err)
	}
	return ln
}

func StartGrpcMockupServer() (context.Context, *grpc.ClientConn) {
	// start GRPC mockup server
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithContextDialer(dialer()))
	if err != nil {
		log.Fatal(err)
	}
	return ctx, conn
}

func dialer() func(context.Context, string) (net.Conn, error) {
	listener := bufconn.Listen(1024 * 1024)
	server := grpc.NewServer()

	var opiSpdkServer 

	pb.RegisterFrontendNvmeServiceServer(server, &opiSpdkServer)
	pb.RegisterMiddleendServiceServer(server, &opiSpdkServer)
	pb.RegisterNVMfRemoteControllerServiceServer(server, &opiSpdkServer)

	go func() {
		if err := server.Serve(listener); err != nil {
			log.Fatal(err)
		}
	}()

	return func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}
}
