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
	filePath := "/default"
	bucket := "ramein-stg"
	fileName := ""

	req, err := conn.Recv()
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	if req.GetFileName() == "" {
		return status.Error(codes.InvalidArgument, "file name must not be empty")
	}

	fileName = req.GetFileName()

	if req.GetFilePath() != "" {
		filePath = req.GetFilePath()
	}

	if req.GetBucket() != "" {
		bucket = req.GetBucket()
	}

	for {
		req, err := conn.Recv()

		if err == io.EOF {

			file, err := os.Create(fmt.Sprintf("%s", fileName))
			//err = s.fileUc.UploadFile(req)
			if err != nil {
				resp := &proto.UploadFileResponse{
					Status:   "Failed",
					FileName: fileName,
					FilePath: filePath,
					Bucket:   bucket,
					Message:  err.Error(),
				}

				return conn.SendAndClose(resp)
			}

			_, err = imageData.WriteTo(file)
			if err != nil {
				resp := &proto.UploadFileResponse{
					Status:   "Failed",
					FileName: fileName,
					FilePath: filePath,
					Bucket:   bucket,
					Message:  err.Error(),
				}

				return conn.SendAndClose(resp)
			}

			resp := &proto.UploadFileResponse{
				Status:   "Success",
				FileName: fileName,
				FilePath: filePath,
				Bucket:   bucket,
				Message:  "Successfully Uploaded file to google cloud storage",
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
