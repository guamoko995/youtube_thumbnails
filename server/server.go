package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"

	"google.golang.org/grpc"

	pb "youtube_thumbnail/thumbnail"
)

var (
	serverAddr = flag.String("addr", "localhost:50051", "The server address in the format of host:port")
	port       = flag.Int("port", 50051, "The server port")
)

type thumbnailServer struct {
	pb.UnimplementedThumbailServer
	//Cache map[string]byte
}

// GetFeature returns the feature at the given point.
func (s *thumbnailServer) GetThumbnail(ctx context.Context, url *pb.Url) (*pb.Img, error) {
	imgUrl := strings.Replace(url.Val, "youtu.be", "img.youtube.com/vi", 1) + "/hqdefault.jpg"

	resp, err := http.Get(imgUrl)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)

	return &pb.Img{Val: buf.Bytes()}, nil

}

func newServer() *thumbnailServer {
	s := &thumbnailServer{} //Cache: make(map[string][]*pb.Msg)}
	return s
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterThumbailServer(grpcServer, newServer())
	grpcServer.Serve(lis)
}
