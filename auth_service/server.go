package main

import (

	pb "app/pb"
	"flag"
  common	"app/common"
	"log"
	"net"
	"google.golang.org/grpc"
)
func main(){

	var (
		environment bool
		listenAddr *string
		jwtPrivateKey   = flag.String("jwt-private-key", "jwtRS256.pem", "The RSA private key to use for signing JWTs")
	)

	if environment = false ; environment == false {
			listenAddr = flag.String("listen-addr", ":7800", "HTTP listen address.")
	} else {
		  listenAddr = flag.String("listen-addr", ":7800", "HTTP listen address.")
	}

	log.Println("Auth service starting...")
	pathConfig := "./Configs.json"
	common.LoadConfiguration(pathConfig,environment)

	gs := grpc.NewServer()

	as, err := NewAuthServer(*jwtPrivateKey)
	if err != nil {
		log.Fatal(err)
	}
	pb.RegisterAuthServer(gs, as)

	ln,err := net.Listen("tcp",*listenAddr)
	if err != nil {
			log.Fatal(err)
	}

	gs.Serve(ln)
}