package main

import (
	"fmt"
	"flag"
	"google.golang.org/grpc"
	"log"
	pb "../pb"
	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"
)



func main(){
	var (
		// caCert     = flag.String("ca-cert", withConfigDir("ca.pem"), "Trusted CA certificate.")
		// serverAddr = flag.String("server-addr", "172.17.0.2:7800", "Auth service address.")
		serverAddr = flag.String("server-addr", "127.0.0.1:7801", "Auth service address.")
		token = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImFiY0B0eW5rLmNvbSIsImV4cCI6MTUyMTU2NTk0MywiaWF0IjoxNTIxMzA2NzQzLCJpc3MiOiJhdXRoLnNlcnZpY2UiLCJzdWIiOiJhYmMxMjMifQ.GdgieeIL-7HPkz9fIdgo2APE1IwtCSruB6Q83DvJoU0foruqv7al9l_PieJnopMGUEzSmFO8qcGe_qzcuIu_uuJaIkaRUkZ7diZxkoFmIw-MDS70cTGFIKalm0QRLJ3QiE7CIZZgmS_MBR_NZkOkXZHW6H01hh5ZSsOJ-xi2K1S9z2krxFhy_NpdQHXIKINAkPFZSC1doS4a9tdAErY3e3z4_fUEvwNxXf1tZOPDAj9z6_VkY4EhwAei00ClGSV_RKdl4N9fzaqzxa1VvzWJ5OI3QGm3jurFtTHhIHTWZP3P_oJlwVrEJ-B3l8PdGq5yUVt07oPUa6omqzDMiEnc4Shk_v5yQ8onEDmvwOM1_4Pg8BCfXEOxIskz7bvoyDFmSn_dRS5BmVm6hv0EALDhAqWfpmQPkQ9wgXyaJ3opX1mFlqK_IWedEi2-6CiacBMNFKHBmoq801M2JjdEgIRbuDaOeF9LLV2Onj6SFbrKC2dcr72oGkCRWo7jhnzKksHVv0PhYsoOradpSUlJX0o01R4KOg-yvrHuoGHmhs_TwXReY6sGM58Pp_bhFVHXrNe7SYbPO97QMjj9L5EfSEuqV4yRt1uJGf7V_hxcxFVwgfTiesaHYR44jb2kn3J4J7uwQg-l5OfizjjXGwzOb2WmWYQ0u0tMjmCbXLkgHaTkDpA"
	)
	flag.Parse()

	// ta, err := credentials.NewClientTLSFromFile(*caCert, "")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	conn, err := grpc.Dial(*serverAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	ac := pb.NewBearClient(conn)
	req := &pb.BearRequest{
		UserId: 3,
	}
	md := metadata.Pairs(
		"Authorization", token,
	)
	ctx := metadata.NewOutgoingContext(context.Background(),md)

	lm,err := ac.GetBearDetail(ctx,req)
	if err != nil {
			log.Fatal(err)
	}

	fmt.Println(lm)
}