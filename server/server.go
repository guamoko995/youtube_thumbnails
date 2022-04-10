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

	db "youtube_thumbnail/database"
	proto "youtube_thumbnail/thumbnail"
)

var (
	serverAddr = flag.String("addr", "localhost:50051", "The server address in the format of host:port")
	port       = flag.Int("port", 50051, "The server port")
	database   = flag.String("database", "thumbnails.sql", "The file of database SQLite")
)

// gRPC прокси-сервис для загрузки thumbnail’ов
// (превью видеоролика) с YouTube.
type thumbnailServer struct {
	proto.UnimplementedThumbailServer
	Cache *db.DB
}

// Реализация единственной RPC сервиса - получения thumbnail’а по url видеоролика Yuotube.
func (s *thumbnailServer) GetThumbnail(ctx context.Context, url *proto.Url) (*proto.Img, error) {
	// Получение thumbnail'а из базы данных (кэш) по ссылке на видеоролик
	img, _ := s.Cache.Get(url.Val)
	if img != nil {
		// В случае успеха возвращает кэшированный thumbnail
		return &proto.Img{Val: img}, nil
	}

	// Если  thumbnail отсутствует в кэше, получает его с Youtub'а
	img, err := getImgFromYoutube(url.Val)
	if err != nil {
		log.Printf("Unable to get image from youtube:   %s/n", err)
		return nil, fmt.Errorf("Ошибка сервера")
	}

	// Кэширует полученный thumbnail
	fmt.Println(s.Cache.Add(url.Val, img))

	return &proto.Img{Val: img}, nil
}

// Получение thumbnail'а с Youtube по ссылке на видеоролик
func getImgFromYoutube(videoUrl string) ([]byte, error) {
	// На серверах Youtube, thumbnail видеоролика
	// c url: https://youtu.be/EXAMPLE
	// имеет url: https://img.youtube.com/vi/EXAMPLE/hqdefault.jpg
	// т.о. url видеоролика трансвформируем в url
	// соответствующего ему thumbnail'а
	thumbnailUrl := strings.Replace(videoUrl, "youtu.be", "img.youtube.com/vi", 1) + "/hqdefault.jpg"

	// Получаем thumbnail по http
	resp, err := http.Get(thumbnailUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Получаем из ответа изображение (буфер байтов)
	img := new(bytes.Buffer)
	img.ReadFrom(resp.Body)

	// Возвращаем изображение в виде массива байтов
	return img.Bytes(), nil
}

// создает экземпляр gRPC прокси-сервиса
func newServer(database string) *thumbnailServer {
	b, err := db.NewDB(database)
	if err != nil {
		panic(err)
	}
	s := &thumbnailServer{Cache: b}
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
	proto.RegisterThumbailServer(grpcServer, newServer(*database))
	grpcServer.Serve(lis)
}
