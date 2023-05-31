package service

import (
	"context"
	"fmt"
	"github.com/achekslime/core/rest_api_utils"
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
	"net/http"
	"net/url"
)

func (srv *MinioService) GetFile(ctx *gin.Context) {
	bucket := ctx.Param("bucket")
	file := ctx.Param("file")

	err := getFile(bucket, file)
	if err != nil {
		rest_api_utils.BindInternalError(ctx, fmt.Errorf("minio err: %s", err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, "success")
}

func getFile(bucket, filename string) error {
	endpoint := "81.200.150.77:9000"
	accessKeyID := "achek"
	secretAccessKey := "qwerty123456"

	s3Client, err := minio.New(
		endpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
			Secure: false,
		})
	if err != nil {
		log.Fatalln(err)
	}

	reqParams := make(url.Values)
	reqParams.Set("response-content-disposition", "attachment; filename=\"minecraft.mp4\"")

	if err := s3Client.FGetObject(context.Background(), bucket, filename, filename, minio.GetObjectOptions{}); err != nil {
		log.Fatalln(err)
	}
	return nil
}
