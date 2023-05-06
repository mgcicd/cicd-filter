package server

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net"
	"time"

	authV3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	"github.com/mgcicd/cicd-core/logs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
)

type AuthService struct {
}

var kaep = keepalive.EnforcementPolicy{
	MinTime:             5 * time.Second, // If a client pings more than once every 5 seconds, terminate the connection
	PermitWithoutStream: true,            // Allow pings even when there are no active streams
}

var kasp = keepalive.ServerParameters{
	MaxConnectionIdle: 180 * time.Second, // If a client is idle for 15 seconds, send a GOAWAY
	Time:              30 * time.Second,  // Ping the client if it is idle for 5 seconds to ensure the connection is still active
	Timeout:           5 * time.Second,   // Wait 1 second for the ping ack before assuming the connection is dead
}

func NewGrpcAuthService(port string) {
	//生成grpc 校验
	ctx := context.Background()

	cert, err := tls.LoadX509KeyPair("certs/server.pem", "certs/server.key")
	if err != nil {
		panic(err)
	}

	certPool := x509.NewCertPool()
	rootBuf, err := ioutil.ReadFile("certs/ca.crt")
	if err != nil {
		panic(err)
	}
	if !certPool.AppendCertsFromPEM(rootBuf) {
		panic("Fail to append ca")
	}

	tlsConf := &tls.Config{
		ClientAuth:   tls.RequireAndVerifyClientCert,
		Certificates: []tls.Certificate{cert},
		ClientCAs:    certPool,
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		logs.DefaultConsoleLog.Error("cicd-filter listen error", fmt.Sprintf("err : %s", err))
	}
	grpcServer := grpc.NewServer(grpc.Creds(credentials.NewTLS(tlsConf)), grpc.KeepaliveEnforcementPolicy(kaep), grpc.KeepaliveParams(kasp))
	server := AuthService{}

	authV3.RegisterAuthorizationServer(grpcServer, server)
	logs.DefaultConsoleLog.Info("cicd-filter", "cicd-filter is listening port 50051")

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			logs.DefaultConsoleLog.Error("cicd-filter", fmt.Sprintf("err : %s", err))
		}
	}()

	<-ctx.Done()

	grpcServer.GracefulStop()
}
