package service

import (
	"bytes"
	"fmt"
	"github.com/andibalo/ramein/corvus/internal/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"os"
)

type FileService struct {
	proto.FileServer
}

func NewFileService() *FileService {
	return &FileService{}
}

func (s *FileService) UploadFile(conn proto.File_UploadFileServer) error {
	imageData := bytes.Buffer{}
	imageSize := 0

	for {
		req, err := conn.Recv()

		if err == io.EOF {

			if req.Bucket == "" {
				req.Bucket = "ramein-stg"
			}

			if req.FilePath == "" {
				req.FilePath = "/default"
			}

			req.File = imageData.Bytes()

			file, err := os.Create(fmt.Sprintf("/%s", req.FileName))
			//err = s.fileUc.UploadFile(req)
			if err != nil {
				resp := &proto.UploadFileResponse{
					Status:   "Failed",
					FileName: req.FileName,
					FilePath: req.FilePath,
					Bucket:   req.Bucket,
					Message:  err.Error(),
				}

				return conn.SendAndClose(resp)
			}

			_, err = imageData.WriteTo(file)
			if err != nil {
				resp := &proto.UploadFileResponse{
					Status:   "Failed",
					FileName: req.FileName,
					FilePath: req.FilePath,
					Bucket:   req.Bucket,
					Message:  err.Error(),
				}

				return conn.SendAndClose(resp)
			}

			resp := &proto.UploadFileResponse{
				Status:   "Success",
				FileName: req.FileName,
				FilePath: req.FilePath,
				Bucket:   req.Bucket,
				Message:  "Successfully Uuplaoded file to google clodu storage",
			}

			return conn.SendAndClose(resp)
		}
		if err != nil {
			return err
		}

		chunk := req.GetFile()
		size := len(chunk)

		imageSize += size
		if imageSize > 1<<20 {
			return status.Errorf(codes.Internal, "image size too large", err)
		}
		_, err = imageData.Write(req.GetFile())
		if err != nil {
			return status.Error(codes.Internal, err.Error())
		}
	}

	return nil
}
