package main

import (
	"context"
	"io"
	"mime/multipart"
)

func uploadData(
	config appConfig, objectId string, f *multipart.FileHeader,
) (int64, error) {
	ctx := context.Background()

	fData, err := f.Open()
	if err != nil {
		return 0, err
	}
	defer fData.Close()

	w, err := config.packageBucket.NewWriter(ctx, objectId, nil)
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
