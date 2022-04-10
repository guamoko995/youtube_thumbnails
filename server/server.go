package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"

	"google.golang.org/grpc"

	pb "youtube_thumbnail/thumbnail"
)

var (
	serverAddr = flag.String("addr", "localhost:50051", "The server address in the format of host:port")
	port       = flag.Int("port", 50051, "The server port")
)

type thumbnailServer struct {
	pb.UnimplementedThumbailServer
	//savedUrl []*pb.Url // read-only after initialized

	mu sync.Mutex // protects routeNotes
	//Cache map[string][]*pb.Url
}

// GetFeature returns the feature at the given point.
func (s *thumbnailServer) GetThumbnail(ctx context.Context, url *pb.Url) (*pb.Img, error) {
	imgUrl := strings.Replace(url.Val, "youtu.be", "img.youtube.com/vi", 1) + "/hqdefault.jpg"

	// https://img.youtube.com/vi/z-mHhobE0Pw/hqdefault.jpg

	fmt.Println(imgUrl)
	//fname += ".jpg"

	resp, err := http.Get(imgUrl)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	out := make([]byte, 0)
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		out = append(out, []byte(scanner.Text())...)
	}
	return &pb.Img{Val: []byte(out)}, nil
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
