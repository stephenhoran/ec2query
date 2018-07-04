package apis

import (
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type s3Bucket struct {
	Name     string
	Location string
	Size     int64
}

var s3slice []s3Bucket

func GetS3Buckets() {
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
		s3slice = append(s3slice, s3Bucket{
			Name:     aws.StringValue(b.Name),
			Location: aws.StringValue(resp.LocationConstraint),
		})
	}

	c := make(chan s3Bucket)

	for _, i := range s3slice {
		go func(i s3Bucket) {
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

			i.Size = size
			c <- i
		}(i)
	}

	for n := 0; n < len(s3slice); n++ {
		fmt.Println(<-c)
	}

}
