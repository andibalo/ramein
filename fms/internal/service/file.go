package service

import (
	"bytes"
	"errors"
	"fms/internal/biz"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"

	v1 "fms/api/file/v1"
)

type FileService struct {
	v1.UnimplementedFileServer
	fileUc *biz.FileUsecase
}

func NewFileService(fileUc *biz.FileUsecase) *FileService {
	return &FileService{
		fileUc: fileUc,
	}
}

func (s *FileService) UploadFile(conn v1.File_UploadFileServer) error {
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

			err = s.fileUc.UploadFile(req)
			if err != nil {
				resp := &v1.UploadFileReply{
					Status:   "Failed",
					FileName: req.FileName,
					FilePath: req.FilePath,
					Bucket:   req.Bucket,
					Message:  err.Error(),
				}

				return conn.SendAndClose(resp)
			}

			resp := &v1.UploadFileReply{
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
			return errors.New("image size too large")
		}
		_, err = imageData.Write(req.GetFile())
		if err != nil {
			return status.Error(codes.Internal, err.Error())
		}
	}
}
