package apis

import (
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// S3Bucket for struct type containing information about S3 buckets
type S3Bucket struct {
	Name     string
	Location string
	Size     string
}

var s3slice []S3Bucket
var s3sliceComplete []S3Bucket

// GetS3Buckets function used to grab a list of S3 buckets and their sizes
func GetS3Buckets() []S3Bucket {
	sess, err := session.NewSession(&aws.Config{Region: aws.String(os.Getenv("REGION"))})
	if err != nil {
		panic(err)
	}

	svc := s3.New(sess)

	result, err := svc.ListBuckets(nil)
	if err != nil {
		panic(err)
	}

	for _, b := range result.Buckets {
		resp, err := svc.GetBucketLocation(&s3.GetBucketLocationInput{Bucket: b.Name})
		if err != nil {
			panic(err)
		}
		s3slice = append(s3slice, S3Bucket{
			Name:     aws.StringValue(b.Name),
			Location: aws.StringValue(resp.LocationConstraint),
		})
	}

	c := make(chan S3Bucket)

	for _, i := range s3slice {
		go func(i S3Bucket) {
			var sess *session.Session
			var size int64
			if i.Location != "" {
				sess, _ = session.NewSession(&aws.Config{Region: aws.String(i.Location)})
			} else {
				sess, _ = session.NewSession(&aws.Config{Region: aws.String(os.Getenv("REGION"))})
			}

			svc := s3.New(sess)
			res, err := svc.ListObjects(&s3.ListObjectsInput{Bucket: &i.Name})
			if err != nil {
				log.Panicf("Unable to list items in bucket %q, %v", i.Name, err)
			}

			for _, item := range res.Contents {
				size += aws.Int64Value(item.Size)
			}
			fsize := float64(size)
			i.Size = fmt.Sprintf("%6.2f MiB", (fsize / 1024 / 1024))
			c <- i
		}(i)
	}

	for n := 0; n < len(s3slice); n++ {
		s3sliceComplete = append(s3sliceComplete, <-c)
	}

	fmt.Println(s3sliceComplete)

	close(c)

	return s3sliceComplete

}
