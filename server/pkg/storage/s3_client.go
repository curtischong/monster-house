package storage

import (
	"github.com/aws/aws-sdk-go/aws"
)

type S3Client struct {
	bucketName string
	config     *aws.Config
}

func NewS3Client(
	config *aws.Config, bucketName string,
) *S3Client {
	return &S3Client{
		bucketName: bucketName,
		config:     config,
	}
}

func (client *S3Client) uploadFile(){

}
