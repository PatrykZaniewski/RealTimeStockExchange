package main

import (
	pb "RealTimeStockExchange/Client/grpc_generated"
	"context"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
	"time"
)

const (
	defaultName = "world"
)

var (
	addr = flag.String("addr", "rtse-server-elej43qtya-lz.a.run.app:443", "the address to connect to")
	name = flag.String("name", defaultName, "Name to greet")
)

func main() {
	flag.Parse()

	systemRoots, err := x509.SystemCertPool()
	if err != nil {
		log.Println("cert")
	}
	creds := credentials.NewTLS(&tls.Config{
		RootCAs: systemRoots,
	})
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	conn, err := grpc.DialContext(ctx, *addr,
		grpc.WithTransportCredentials(creds),
	)

	defer conn.Close()
	if err != nil {
		log.Println("AA")
	}
	c := pb.NewGreeterClient(conn)

	defer cancel()
	r, err := c.SayHello(ctx, &pb.Request{Request: *name})
	if err != nil {
		log.Fatalf("could not greet: %s", err)
	}
	log.Printf("Greeting: %s", r.Response)
}
