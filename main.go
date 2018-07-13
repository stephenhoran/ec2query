package main

import (
	"context"
	"os"
	"runtime"

	"github.com/atssteve/ec2query/apis"
	"github.com/atssteve/ec2query/email"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// Handle used by the lambda
func Handle() {

	ctx := context.Background()

	runtime.GOMAXPROCS(runtime.NumCPU())

	data := make(map[string]interface{})

	sess, err := session.NewSession(&aws.Config{Region: aws.String(os.Getenv("REGION"))})
	if err != nil {
		panic(err)
	}

	// Starting a new EC2 client.
	ec2session := ec2.New(sess)

	// // Describe all of the regions in AWS. Returns a type *DescribeRegionsOutput.
	resultRegions := apis.GetRegions(ctx, ec2session)

	// // Query all regions for all instance types and and recieve a slice of maps for that region.
	data["instances"] = apis.GetInstances(resultRegions)
	data["s3"] = apis.GetS3Buckets()

	// // Build the HTML body of the email.
	htmlBody := email.BuildInstanceTemplate(data)

	// Send Email
	apis.SendEmail(sess, htmlBody)

}

func main() {
	lambda.Start(Handle)
}
