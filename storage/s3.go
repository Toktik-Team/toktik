package storage

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"io"
	"log"
	"net/url"
	"strconv"
	appcfg "toktik/constant/config"
)

var client *s3.Client

func init() {
	r2Resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL: appcfg.EnvConfig.S3_ENDPOINT_URL,
		}, nil
	})

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithEndpointResolverWithOptions(r2Resolver),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(appcfg.EnvConfig.S3_SECRET_ID, appcfg.EnvConfig.S3_SECRET_KEY, "")),
	)
	if err != nil {
		log.Fatal(err)
	}

	client = s3.NewFromConfig(cfg, func(o *s3.Options) {
		// Required when using minio
		o.UsePathStyle, _ = strconv.ParseBool(appcfg.EnvConfig.S3_PATH_STYLE)
	})
}

// Upload to the s3 storage using given fileName
func Upload(fileName string, content io.Reader) (*s3.PutObjectOutput, error) {
	resp, err := client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: &appcfg.EnvConfig.S3_BUCKET,
		Key:    &fileName,
		Body:   content,
	})
	return resp, err
}

func GetLink(fileName string) (string, error) {
	return url.JoinPath(appcfg.EnvConfig.S3_PUBLIC_URL, fileName)
}
