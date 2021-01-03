package storage

import (
	"../config"
	"fmt"
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

func (client *S3Client) GetAllFileURLs()([]string, error){
	s3client := client.getS3Client()
	// 1) Get all of the objects in our bucket
	result, err := s3client.ListObjects(&s3.ListObjectsInput{
		Bucket: aws.String(client.awsConfig.S3BucketName),
	})
	if err != nil{
		return nil, err
	}
	fileUrls := make([]string, len(result.Contents))
	for i, fileContent := range result.Contents{
		if fileContent.Key == nil{
			continue
		}
		fileUrl := fmt.Sprintf("%s/%s/%s",client.awsConfig.S3Endpoint,
			client.awsConfig.S3BucketName, *fileContent.Key)
		fileUrls[i] = fileUrl
	}
	return fileUrls, nil
}

func (client *S3Client) UploadFile(
	f io.ReadSeeker, fileType string,
)(fileID uuid.UUID, err error){
	s3client := client.getS3Client()
	fileID = uuid.New()
	_, err = s3client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(client.awsConfig.S3BucketName),
		ContentType: aws.String(fileType),
		ACL: aws.String("public-read"),
		Key:    aws.String(fileID.String()),
		Body:   f,
	})
	return
}

func (client *S3Client) getS3Client()*s3.S3{
	awsConfig := &aws.Config{
		Region: aws.String(client.awsConfig.Region),
		Endpoint: aws.String(client.awsConfig.S3Endpoint),
		S3ForcePathStyle: aws.Bool(true),
	}
	return s3.New(session.Must(session.NewSession(awsConfig)))
}
