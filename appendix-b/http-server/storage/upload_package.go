package storage

import (
	"context"
	"io"
	"mime/multipart"

	"github.com/practicalgo/code/appendix-b/http-server/config"
)

func UploadData(
	config *config.AppConfig, objectId string, f *multipart.FileHeader,
) (int64, error) {
	ctx := context.Background()

	fData, err := f.Open()
	if err != nil {
		return 0, err
	}
	defer fData.Close()

	w, err := config.PackageBucket.NewWriter(ctx, objectId, nil)
	if err != nil {
		return 0, err
	}

	nBytes, err := io.Copy(w, fData)
	if err != nil {
		return 0, err
	}
	err = w.Close()
	if err != nil {
		return nBytes, err
	}
	return nBytes, nil
}
