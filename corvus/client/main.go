package main

import (
	"bufio"
	"context"
	pb "github.com/andibalo/ramein/corvus/internal/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

func main() {

	cl, err := grpc.Dial("0.0.0.0:8082", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		panic("error connecting to grpc server")
	}

	fileCl := NewFileClient(cl)

	fileCl.UploadImage("tmp/test.jpg")
}

type FileClient struct {
	service pb.FileClient
}

// NewFileClient returns a new file client
func NewFileClient(cc *grpc.ClientConn) *FileClient {
	service := pb.NewFileClient(cc)
	return &FileClient{service}
}

// UploadImage calls upload image RPC
func (c *FileClient) UploadImage(imagePath string) {
	file, err := os.Open(imagePath)
	if err != nil {
		log.Fatal("cannot open image file: ", err)
	}
	defer file.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	stream, err := c.service.UploadFile(ctx)
	if err != nil {
		log.Fatal("cannot upload image: ", err)
	}

	req := &pb.UploadFileRequest{
		FileName: filepath.Base(imagePath),
	}

	err = stream.Send(req)
	if err != nil {
		log.Fatal("cannot send image info to server: ", err, stream.RecvMsg(nil))
	}

	reader := bufio.NewReader(file)
	buffer := make([]byte, 1024)

	for {
		n, err := reader.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal("cannot read chunk to buffer: ", err)
		}

		req := &pb.UploadFileRequest{
			File: buffer[:n],
		}

		err = stream.Send(req)
		if err != nil {
			log.Fatal("cannot send chunk to server: ", err, stream.RecvMsg(nil))
		}
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatal("cannot receive response: ", err)
	}

	log.Printf("%v", res)
}
