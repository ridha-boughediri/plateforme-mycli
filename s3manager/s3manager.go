package s3manager

import (
    "context"
    "fmt"
    "os"

    "github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/credentials"
    "github.com/aws/aws-sdk-go-v2/service/s3"
)

func CreateBucket(bucketName string) error {
    cfg, err := config.LoadDefaultConfig(context.TODO(),
        config.WithRegion(os.Getenv("AWS_REGION")),
        config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
            os.Getenv("AWS_ACCESS_KEY_ID"), os.Getenv("AWS_SECRET_ACCESS_KEY"), "")),
        config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
            func(service, region string, options ...interface{}) (aws.Endpoint, error) {
                if service == s3.ServiceID {
                    return aws.Endpoint{
                        URL:           os.Getenv("S3_ENDPOINT"),
                        SigningRegion: os.Getenv("AWS_REGION"),
                        HostnameImmutable: true,
                    }, nil
                }
                return aws.Endpoint{}, &aws.EndpointNotFoundError{}
            })),
    )
    if err != nil {
        return fmt.Errorf("unable to load SDK config, %v", err)
    }

    s3Client := s3.NewFromConfig(cfg)
    _, err = s3Client.CreateBucket(context.TODO(), &s3.CreateBucketInput{
        Bucket: &bucketName,
    })

    return err
}





// ListBuckets lists all S3 buckets with detailed request/response logging
func ListBuckets() error {
    // Load AWS configuration with detailed logging enabled
    cfg, err := config.LoadDefaultConfig(context.TODO(),
        config.WithRegion(os.Getenv("AWS_REGION")),
        config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
            os.Getenv("AWS_ACCESS_KEY_ID"), os.Getenv("AWS_SECRET_ACCESS_KEY"), "")),
        config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
            func(service, region string, options ...interface{}) (aws.Endpoint, error) {
                if service == s3.ServiceID {
                    return aws.Endpoint{
                        URL:           os.Getenv("S3_ENDPOINT"),
                        SigningRegion: os.Getenv("AWS_REGION"),
                    }, nil
                }
                return aws.Endpoint{}, &aws.EndpointNotFoundError{}
            })),
        config.WithClientLogMode(aws.LogRequestWithBody | aws.LogResponseWithBody), // Enable detailed logging
    )
    if err != nil {
        return fmt.Errorf("unable to load SDK config, %v", err)
    }

    // Create S3 client
    s3Client := s3.NewFromConfig(cfg)

    // Call ListBuckets
    result, err := s3Client.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
    if err != nil {
        return fmt.Errorf("unable to list buckets, %v", err)
    }

    // Output the list of buckets
    for _, bucket := range result.Buckets {
        fmt.Println(*bucket.Name)
    }

    return nil
}
