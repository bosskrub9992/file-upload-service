package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/bosskrub9992/file-upload-service/objectstores"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

func main() {
	ctx := context.Background()
	storageClient, err := storage.NewClient(ctx, option.WithCredentialsFile("./key.json"))
	if err != nil {
		slog.Error(err.Error())
		return
	}
	defer storageClient.Close()
	gcs := objectstores.NewGoogleCloudStorage(storageClient)
	f, err := os.Open("./cmd/testgcs/please.jpeg")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	srcByte, err := io.ReadAll(f)
	if err != nil {
		fmt.Println(err)
		return
	}
	if err := gcs.Upload(ctx, srcByte, "public_bucket_123", "please.jpeg", "image/jpeg"); err != nil {
		fmt.Println(err)
		return
	}
}
