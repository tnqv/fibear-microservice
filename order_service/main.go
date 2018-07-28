package main

import(
	"flag"
	common "./common"
	"log"
	"google.golang.org/grpc"
	pb "./pb"
	"net"
)

func main(){
	var (
		environment bool
		listenAddr *string
		jwtPublicKey   = flag.String("jwt-public-key", "jwtRS256.key.pub", "The RSA public key to use for signing JWTs")
	)

	if environment = false ; environment == false {
			listenAddr = flag.String("listen-addr", "127.0.0.1:7802", "HTTP listen address.")
	} else {
			listenAddr = flag.String("listen-addr", "0.0.0.0:7802", "HTTP listen address.")
	}

	log.Println("Order service starting...")
	pathConfig := "./Configs.json"
	common.LoadConfiguration(pathConfig,environment)

	gs := grpc.NewServer()

	bs,err := NewOrderServer(*jwtPublicKey)
	if err != nil {
			log.Fatal(err)
	}

	pb.RegisterOrderServer(gs,bs)

	ln,err := net.Listen("tcp",*listenAddr)
	if err != nil {
			log.Fatal(err)
	}

	gs.Serve(ln)


}