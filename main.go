package main

import (
	"cicd-filter/server"
	"flag"
	_ "net/http/pprof"
)

var GrpcPortFlag = flag.String("authGrpcPort", "50052", "authServer-gprc-Port")

func main() {

	server.NewGrpcAuthService(*GrpcPortFlag)

}
