package main

import (
	"runtime"

	"github.com/atssteve/ec2query/apis"
)

// Handle used by the lambda
func Handle() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	// sess, err := session.NewSession(&aws.Config{Region: aws.String(os.Getenv("REGION"))})
	// if err != nil {
	// 	panic(err)
	// }

	// // Starting a new EC2 client.
	// ec2session := ec2.New(sess)

	// // Describe all of the regions in AWS. Returns a type *DescribeRegionsOutput.
	// resultRegions := apis.GetRegions(ec2session)

	// // Query all regions for all instance types and and recieve a slice of maps for that region.
	// resultInstances := apis.GetInstances(resultRegions)

	// // Build the HTML body of the email.
	// htmlBody := email.BuildInstanceTemplate(resultInstances)

	// apis.SendEmail(sess, htmlBody)
	apis.GetS3Buckets()
}

func main() {
	Handle()
}
