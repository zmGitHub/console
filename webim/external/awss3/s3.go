package awss3

import (
	"fmt"
	"io"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"bitbucket.org/forfd/custm-chat/webim/common/logging"
	"bitbucket.org/forfd/custm-chat/webim/conf"
)

const (
	S3Host   = `https://chat-im.s3.ap-southeast-1.amazonaws.com`
	FilePath = S3Host + "/%s/files/" // /{ent_id}/files/
)

type UploadClient struct {
	region          string
	bucketName      string
	accessKeyID     string
	secretAccessKey string
	sessionToken    string
}

func NewUploadClient(config *conf.AWSS3Config) (cli *UploadClient) {
	cli = &UploadClient{
		region:          config.Region,
		bucketName:      config.BucketName,
		accessKeyID:     config.AccessKeyID,
		secretAccessKey: config.SecretAccessKey,
		sessionToken:    "",
	}
	return
}

func (cli *UploadClient) Upload(fileName string, file io.Reader, expireAt time.Time) (string, error) {
	sess := session.Must(session.NewSession(&aws.Config{
		Region:      aws.String(cli.region),
		Credentials: credentials.NewStaticCredentials(cli.accessKeyID, cli.secretAccessKey, cli.sessionToken),
		MaxRetries:  aws.Int(3),
		Logger: aws.LoggerFunc(func(args ...interface{}) {
			log.Logger.Info(args...)
		}),
	}))

	uploader := s3manager.NewUploader(sess)
	input := &s3manager.UploadInput{
		Bucket: aws.String(cli.bucketName),
		Key:    aws.String(fileName),
		Body:   file,
	}
	if !expireAt.IsZero() {
		input.Expires = &expireAt
	}

	result, err := uploader.Upload(input)
	if err != nil {
		return "", err
	}

	if result == nil {
		return "", fmt.Errorf("upload result nil")
	}

	return result.Location, nil
}
