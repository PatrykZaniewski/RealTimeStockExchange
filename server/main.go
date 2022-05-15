package main

import (
	pb "RealTimeStockExchange/Server/grpc_generated"
	"context"
	"flag"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"log"
	"net"
	"net/http"
	"os"
	"time"
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

func getAlbums(c *gin.Context) {
	c.IndentedJSON(http.StatusCreated, "AB")
}

func main() {
	//router := gin.Default()
	//router.GET("/albums", getAlbums)
	//
	//router.Run("0.0.0.0:" + os.Getenv("PORT_RUN"))

	flag.Parse()
	lis, err := net.Listen("tcp", "0.0.0.0:"+os.Getenv("PORT_RUN"))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer(grpc.KeepaliveParams(keepalive.ServerParameters{
		MaxConnectionIdle: 5 * time.Minute,
	}))
	pb.RegisterGreeterServer(s, &GreeterServer{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
