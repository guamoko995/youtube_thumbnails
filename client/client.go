package main

import (
	"context"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	pb "youtube_thumbnail/thumbnail"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	async      = flag.Bool("async", false, "Files are loaded asynchronously if true, otherwise in order")
	out        = flag.String("out", "", "download folder")
	serverAddr = flag.String("addr", "localhost:50051", "The server address in the format of host:port")
)

func GetThumbnail(client pb.ThumbailClient, url string, path string, fname string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	img, err := client.GetThumbnail(ctx, &pb.Url{Val: url})
	if err != nil {
		log.Fatalf("%v.GetThumbnail(_) = _, %v: ", client, err)
	}
	if err := ioutil.WriteFile(path+"\\"+fname, img.Val, 0644); err != nil {
		log.Fatalln("Failed to write file:", err)
	}
}

func main() {
	flag.Parse()
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := pb.NewThumbailClient(conn)

	// Получение основных параметров (ссылок на видео youtube)

	urls := make([]string, 0)
	for _, arg := range flag.Args() {
		urls = append(urls, arg)
	}

	if len(urls) == 0 {
		log.Fatalln("Not set url(s)")
	}

	// Если каталог загрузки не указан, используется текущий
	if *out == "" {
		folder, _ := os.Getwd()
		*out = folder
	}

	// Скачивание thumbnail'ов
	if *async {
		//асинхронное скачивание
		var wg sync.WaitGroup

		// функция асинхронного вызов GetThumbnail с контролем ожидания
		assignGetThumbnail := func(url, out, fname string) {
			wg.Add(1)
			go func() {
				defer func() { wg.Done() }()
				GetThumbnail(client, url, out, fname)
			}()
		}
		for _, url := range urls {
			fname := strings.Replace(url, "https://youtu.be/", "", 1) + ".jpg"
			assignGetThumbnail(url, *out, fname)
		}
		wg.Wait()
	} else {
		//последовательное скачивание
		for _, url := range urls {
			fname := strings.Replace(url, "https://youtu.be/", "", 1)
			GetThumbnail(client, url, *out, fname)
		}
	}
}
