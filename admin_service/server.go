package main

import (

	pb "./pb"
	"flag"
  common	"./common"
	"log"
	"net"
	"google.golang.org/grpc"
)
func main(){

	var (
		environment bool
		listenAddr *string
		jwtPublicKey   = flag.String("jwt-public-key", "jwtRS256.key.pub", "The RSA private key to use for signing JWTs")
		jwtPrivateKey   = flag.String("jwt-private-key", "jwtRS256.pem", "The RSA private key to use for signing JWTs")
	)

	if environment = false ; environment == false {
			listenAddr = flag.String("listen-addr", "127.0.0.1:7807", "HTTP listen address.")
	} else {
		  listenAddr = flag.String("listen-addr", "0.0.0.0:7807", "HTTP listen address.")
	}

	log.Println("Admin service starting...")
	pathConfig := "./Configs.json"
	common.LoadConfiguration(pathConfig,environment)

	gs := grpc.NewServer()

	as, err := NewAdminServer(*jwtPublicKey,*jwtPrivateKey)
	if err != nil {
		log.Fatal(err)
	}
	pb.RegisterAdminServer(gs, as)

	ln,err := net.Listen("tcp",*listenAddr)
	if err != nil {
			log.Fatal(err)
	}

	gs.Serve(ln)
}