package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var (
	endpoint, accessKey, secretKey string
	bucket, prefix                 string
	insecure, recursive            bool
	skipErr                        bool
)

func main() {
	flag.StringVar(&endpoint, "endpoint", "", "S3 endpoint URL")
	flag.StringVar(&accessKey, "access-key", "", "S3 Access Key")
	flag.StringVar(&secretKey, "secret-key", "", "S3 Secret Key")
	flag.StringVar(&bucket, "bucket", "", "Select a specific bucket")
	flag.StringVar(&prefix, "prefix", "", "Select an object/prefix")
	flag.BoolVar(&insecure, "insecure", false, "Disable TLS verification")
	flag.BoolVar(&recursive, "recursive", false, "Enable recursive listing")
	flag.BoolVar(&skipErr, "skiperror", false, "Skip other errors")
	flag.Parse()

	if endpoint == "" {
		log.Fatalln("endpoint is not provided")
	}

	if accessKey == "" {
		log.Fatalln("access key is not provided")
	}

	if secretKey == "" {
		log.Fatalln("secret key is not provided")
	}

	if bucket == "" {
		log.Fatalln("bucket should not be empty")
	}

	s3Client := getS3Client(endpoint, accessKey, secretKey, insecure)

	ctx := context.Background()

	for obj := range s3Client.ListObjects(ctx, bucket, minio.ListObjectsOptions{
		Recursive:    recursive,
		Prefix:       strings.TrimPrefix(prefix, "/"),
		WithVersions: true,
		WithMetadata: true,
	}) {
		if obj.Err != nil {
			log.Fatalln("unable to list with error:", obj.Err)
			return
		}
		if strings.HasSuffix(obj.Key, "/") {
			// directory marker
			continue
		}
		if _, err := s3Client.StatObject(ctx, bucket, obj.Key, minio.StatObjectOptions{}); err != nil {
			if minio.ToErrorResponse(err).Code == "SlowDownRead" && minio.ToErrorResponse(err).StatusCode == 503 {
				fmt.Printf("/%s\n", obj.Key)
				continue
			}
			if skipErr {
				continue
			}
			log.Fatalln(err)
		}
	}
}

func getS3Client(endpoint string, accessKey string, secretKey string, insecure bool) *minio.Client {
	u, err := url.Parse(endpoint)
	if err != nil {
		log.Fatalln(err)
	}

	secure := strings.EqualFold(u.Scheme, "https")
	transport, err := minio.DefaultTransport(secure)
	if err != nil {
		log.Fatalln(err)
	}
	if transport.TLSClientConfig != nil {
		transport.TLSClientConfig.InsecureSkipVerify = insecure
	}

	s3Client, err := minio.New(u.Host, &minio.Options{
		Creds:     credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure:    secure,
		Transport: transport,
	})
	if err != nil {
		log.Fatalln(err)
	}
	return s3Client
}
