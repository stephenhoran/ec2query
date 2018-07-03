package apis

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// Ec2instance struct to hold all of our ec2 instance metadata
type Ec2instance struct {
	Instanceid string
	Type       string
	LaunchTime *time.Time
	State      string
	KeyName    string
	Region     string
}

// GetInstances takes a pointer to a DescribeRegionsOutput as will as a pointer to a string
// to query all of the instances in a region and build on our Ec2instance struct
func GetInstances(regions *ec2.DescribeRegionsOutput, instances *[]Ec2instance) {
	// Iterate over our list of regions and use aws.StringValue to print the region name.
	c := make(chan Ec2instance)
	defer close(c)
	for _, region := range regions.Regions {
		go func() {
			fmt.Println(aws.StringValue(region.RegionName))
			queryInstances(aws.StringValue(region.RegionName), c)
		}()
	}
	go func() {
		for n := range c {
			fmt.Println(n)
		}
	}()
}

// GetInstances returns a list of Ec2instance structs that are currently running
func queryInstances(regionName string, c chan<- Ec2instance) {
	fmt.Println(regionName)
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(regionName),
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
			isstruct := Ec2instance{
				Instanceid: aws.StringValue(instances.InstanceId),
				Type:       aws.StringValue(instances.InstanceType),
				LaunchTime: instances.LaunchTime,
				State:      aws.StringValue(instances.State.Name),
				KeyName:    aws.StringValue(instances.KeyName),
				Region:     aws.StringValue(&regionName),
			}
			c <- isstruct
		}
	}
}
