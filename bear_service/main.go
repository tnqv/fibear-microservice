package main

import(
	"flag"
	common "app/common"
	"log"
	"google.golang.org/grpc"
	pb "app/pb"
	"net"
)

func main(){
	var (
		environment bool
		listenAddr *string
		jwtPublicKey   = flag.String("jwt-public-key", "jwtRS256.key.pub", "The RSA public key to use for signing JWTs")
	)

	if environment = false ; environment == false {
			listenAddr = flag.String("listen-addr", ":7801", "HTTP listen address.")
	} else {
		  listenAddr = flag.String("listen-addr", ":7801", "HTTP listen address.")
	}

	log.Println("Bear service starting...")
	pathConfig := "./Configs.json"
	common.LoadConfiguration(pathConfig,environment)

	gs := grpc.NewServer()

	bs,err := NewBearServer(*jwtPublicKey)
	if err != nil {
			log.Fatal(err)
	}

	pb.RegisterBearServer(gs,bs)

	ln,err := net.Listen("tcp",*listenAddr)
	if err != nil {
			log.Fatal(err)
	}

	gs.Serve(ln)


}