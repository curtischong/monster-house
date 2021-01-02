package storage

import (
	"../config"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
	"io"
)

type S3Client struct {
	awsConfig config.AWSConfig
}

func NewS3Client(
	config *config.Config,
) *S3Client {
	return &S3Client{
		awsConfig: config.AWSConfig,
	}
}

func (client *S3Client) GetAllFileURLs(){

}

func (client *S3Client) GetNewSession()*session.Session{
	awsConfig := &aws.Config{
		Region: aws.String(client.awsConfig.Region),
		Endpoint: aws.String(client.awsConfig.S3Endpoint),
		S3ForcePathStyle: aws.Bool(true),
	}
	return session.Must(session.NewSession(awsConfig))
}

func (client *S3Client) UploadFile(
	f io.ReadSeeker,
)error{
	sess := client.GetNewSession()
	// Create an uploader with the session and default options
	s3client := s3.New(sess)
	result, err := s3client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(client.awsConfig.S3BucketName),
		Key:    aws.String(uuid.New().String()),
		Body:   f,
	})
	if err != nil{
		return err
	}
	print(result)
	return nil
}
