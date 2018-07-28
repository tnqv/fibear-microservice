package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	pb "./pb"

	"golang.org/x/crypto/ssh/terminal"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func withConfigDir(path string) string {
	return filepath.Join(os.Getenv("HOME"), path)
}

func main() {
	var (
		// caCert     = flag.String("ca-cert", withConfigDir("ca.pem"), "Trusted CA certificate.")
		// serverAddr = flag.String("server-addr", "172.17.0.2:7800", "Auth service address.")
		serverAddr = flag.String("server-addr", "127.0.0.1:7800", "Auth service address.")
		username   = flag.String("username", "abc123", "Username to use.")
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

	ac := pb.NewAuthClient(conn)

	fmt.Println("enter password:")
	password, err := terminal.ReadPassword(0)
	if err != nil {
		log.Fatal(err)
	}


	req := &pb.LoginRequest{
		Username: *username,
		Password: string(password),
	}
	lm, err := ac.Login(context.Background(), req)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(lm.User)
	err = ioutil.WriteFile(withConfigDir(".token"), []byte(lm.Token), 0600)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("wrote", withConfigDir(".token"))
}