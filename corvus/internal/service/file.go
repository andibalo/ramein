package service

import (
	"bytes"
	"github.com/andibalo/ramein/corvus/internal/config"
	"github.com/andibalo/ramein/corvus/internal/external"
	"github.com/andibalo/ramein/corvus/internal/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
)

type FileService struct {
	proto.FileServer
	gcs external.GCSRepo
	cfg config.Config
}

func NewFileService(gcs external.GCSRepo, cfg config.Config) *FileService {
	return &FileService{
		gcs: gcs,
		cfg: cfg,
	}
}

func (s *FileService) UploadFile(conn proto.File_UploadFileServer) error {
	imageData := bytes.Buffer{}
	imageSize := 0
	filePath := "default/"
	bucket := "ramein-stg"
	fileName := ""

	req, err := conn.Recv()
	if err != nil {
		s.cfg.Logger().Error("Unable to receive stream", zap.Error(err))
		return status.Error(codes.Internal, err.Error())
	}

	if req.GetFileName() == "" {
		s.cfg.Logger().Error("file name must not be empty")
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
		req, err = conn.Recv()

		if err == io.EOF {

			err = s.gcs.Upload(imageData.Bytes(), bucket, fileName, filePath)
			if err != nil {
				s.cfg.Logger().Error("Unable to upload to google cloud storage", zap.Error(err))

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
			s.cfg.Logger().Error("Unable to receive stream", zap.Error(err))
			return status.Error(codes.Internal, err.Error())
		}

		chunk := req.GetFile()
		size := len(chunk)

		imageSize += size
		if imageSize > 1<<20 {
			s.cfg.Logger().Error("image size too large")
			return status.Errorf(codes.Internal, "image size too large")
		}
		_, err = imageData.Write(req.GetFile())
		if err != nil {
			s.cfg.Logger().Error("Unable to write received file to image buffer", zap.Error(err))
			return status.Error(codes.Internal, err.Error())
		}
	}

	return nil
}
