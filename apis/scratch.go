package apis

import (
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func main() {
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

		if resp.LocationConstraint != nil {
			*sess.Config.Region = aws.StringValue(resp.LocationConstraint)
		}
		res, err := svc.ListObjects(&s3.ListObjectsInput{Bucket: b.Name})
		if err != nil {
			log.Panicf("Unable to list items in bucket %q, %v", b.Name, err)
		}

		for _, item := range res.Contents {
			fmt.Println("Name:         ", *item.Key)
			fmt.Println("Last modified:", *item.LastModified)
			fmt.Println("Size:         ", *item.Size)
			fmt.Println("Storage class:", *item.StorageClass)
			fmt.Println("")
		}
	}

	// bucket := "my-bucket"
	// region, err := s3manager.GetBucketRegion(ctx, sess, bucket, "us-west-2")
	// if err != nil {
	// 	if aerr, ok := err.(awserr.Error); ok && aerr.Code() == "NotFound" {
	// 		fmt.Fprintf(os.Stderr, "unable to find bucket %s's region not found\n", bucket)
	// 	}
	// 	return err
	// }
	// fmt.Printf("Bucket %s is in %s region\n", bucket, region)
}
