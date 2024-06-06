package services

import (
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func getPresignedURL(bucketName, objectKey string) (string, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:                        aws.String(os.Getenv("AWS_REGION")),
		CredentialsChainVerboseErrors: aws.Bool(true),
		S3ForcePathStyle:              aws.Bool(true),
		//SignatureVersion: aws.String("v4"),
	})
	if err != nil {
		return "", err
	}

	svc := s3.New(sess)
	expiration := (10 * time.Minute)

	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	})

	url, err := req.Presign(expiration)
	if err != nil {
		return "", err
	}

	return url, nil
}
