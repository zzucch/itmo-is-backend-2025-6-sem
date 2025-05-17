package s3

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Client struct {
	client     *s3.Client
	bucketName string
}

func NewS3Client() (*Client, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(cfg)

	if _, err := client.ListBuckets(
		context.TODO(),
		&s3.ListBucketsInput{},
	); err != nil {
		log.Fatal(err)
	}

	return &Client{
		client:     client,
		bucketName: "idk-second",
	}, nil
}

func (s *Client) UploadFile(
	ctx context.Context,
	file multipart.File,
	fileName string,
) (string, error) {
	ext := filepath.Ext(fileName)
	newFileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		return "", err
	}

	if _, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(newFileName),
		Body:   bytes.NewReader(buf.Bytes()),
	}); err != nil {
		return "", err
	}

	presignClient := s3.NewPresignClient(s.client)
	presignedReq, err := presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(newFileName),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = 7 * 24 * time.Hour
	})
	if err != nil {
		return "", fmt.Errorf("failed to generate pre-signed URL: %v", err)
	}

	return presignedReq.URL, nil
}
