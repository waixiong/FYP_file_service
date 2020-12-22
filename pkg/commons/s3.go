package commons

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func InitS3Client(ctx context.Context) *s3.S3 {
	backblazeId := os.Getenv("BACKBLAZE_ID")
	backblazeKey := os.Getenv("BACKBLAZE_KEY")
	backblazeS3Endpoint := os.Getenv("BACKBLAZE_S3ENDPOINT")
	backblazeS3Region := os.Getenv("BACKBLAZE_S3REGION")
	// fmt.Println(len(backblazeId))
	if backblazeId == "" {
		backblazeId = "002735e64043a570000000002"
	}
	// fmt.Println(len(backblazeKey))
	if backblazeKey == "" {
		backblazeKey = "K002eRILylOUtX/cm2WjSIxsqbKGWW8"
	}
	// fmt.Println(len(backblazeS3Endpoint))
	if backblazeS3Endpoint == "" {
		backblazeS3Endpoint = "https://s3.us-west-002.backblazeb2.com"
	}
	// fmt.Println(len(backblazeS3Region))
	if backblazeS3Region == "" {
		backblazeS3Region = "us-west-002"
	}
	// endpoint := aws.String(backblazeS3Endpoint)

	s3Config := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(backblazeId, backblazeKey, ""),
		Endpoint:         aws.String(backblazeS3Endpoint),
		Region:           aws.String(backblazeS3Region),
		S3ForcePathStyle: aws.Bool(true),
	}
	newSession := session.New(s3Config)

	s3Client := s3.New(newSession)

	return s3Client
}

func InitCustomS3Client(ctx context.Context, backblazeId string, backblazeKey string) *s3.S3 {
	// backblazeId := os.Getenv("BACKBLAZE_ID")
	// backblazeKey := os.Getenv("BACKBLAZE_KEY")
	backblazeS3Endpoint := os.Getenv("BACKBLAZE_S3ENDPOINT")
	backblazeS3Region := os.Getenv("BACKBLAZE_S3REGION")
	// fmt.Println(len(backblazeId))
	if backblazeId == "" {
		backblazeId = "002735e64043a570000000002"
	}
	// fmt.Println(len(backblazeKey))
	if backblazeKey == "" {
		backblazeKey = "K002eRILylOUtX/cm2WjSIxsqbKGWW8"
	}
	// fmt.Println(len(backblazeS3Endpoint))
	if backblazeS3Endpoint == "" {
		backblazeS3Endpoint = "https://s3.us-west-002.backblazeb2.com"
	}
	// fmt.Println(len(backblazeS3Region))
	if backblazeS3Region == "" {
		backblazeS3Region = "us-west-002"
	}
	// endpoint := aws.String(backblazeS3Endpoint)

	s3Config := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(backblazeId, backblazeKey, ""),
		Endpoint:         aws.String(backblazeS3Endpoint),
		Region:           aws.String(backblazeS3Region),
		S3ForcePathStyle: aws.Bool(true),
	}
	newSession := session.New(s3Config)

	s3Client := s3.New(newSession)

	return s3Client
}
