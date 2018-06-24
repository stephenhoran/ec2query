package main

import (
	"fmt"
	"runtime"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	sess, err := session.NewSession(&aws.Config{Region: aws.String("us-east-2")})
	if err != nil {
		panic(err)
	}

	// Starting a new EC2 client.
	ec2session := ec2.New(sess)

	// Describe all of the regions in AWS. Returns a type *DescribeRegionsOutput.
	resultRegions, err := ec2session.DescribeRegions(nil)
	if err != nil {
		panic(err)
	}

	artifacts := make(map[string][]string)

	// Grab the length of the slice of regions and create a WaitGroup for this.
	// Iterate over our list of regions and use aws.StringValue to print the region name.
	for _, region := range resultRegions.Regions {
		var is []string
		fmt.Println(aws.StringValue(region.RegionName))
		is = getInstances(*region.RegionName)

		artifacts[aws.StringValue(region.RegionName)] = is
	}
	fmt.Println(artifacts)

}

func getInstances(region string) []string {
	var is []string
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		panic(err)
	}
	ec2svc := ec2.New(sess)
	result, err := ec2svc.DescribeInstances(nil)
	if err != nil {
		panic(err)
	}
	for _, reserv := range result.Reservations {
		for _, instances := range reserv.Instances {
			is = append(is, aws.StringValue(instances.InstanceId))
		}
	}

	return is
}
