package main

import (
	pb "RealTimeStockExchange/Server/grpc_generated"
	"context"
	"flag"
	"google.golang.org/grpc"
	_ "google.golang.org/grpc"
	"log"
	"net"
	_ "net"
	"os"
)

var port = flag.Int("port", 50051, "The server port")

type GreeterServer struct {
	pb.UnimplementedGreeterServer
}

func (s *GreeterServer) SayHello(ctx context.Context, in *pb.Request) (*pb.Response, error) {
	log.Printf("Received")
	msg := "Hello" + in.Request
	return &pb.Response{Response: msg}, nil
}

//func getAlbums(c *gin.Context) {
//	c.IndentedJSON(http.StatusCreated, "AB")
//}

func main() {
	//router := gin.Default()
	//router.GET("/albums", getAlbums)

	flag.Parse()
	lis, err := net.Listen("tcp", "0.0.0.0:"+os.Getenv("PORT"))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &GreeterServer{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
