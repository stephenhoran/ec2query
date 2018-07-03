package apis

import (
	"fmt"
	"sync"
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

var wg sync.WaitGroup
var mutex sync.Mutex
var instancesslice *[]Ec2instance

// GetInstances takes a pointer to a DescribeRegionsOutput as will as a pointer to a string
// to query all of the instances in a region and build on our Ec2instance struct
func GetInstances(regions *ec2.DescribeRegionsOutput, instance *[]Ec2instance) {
	// Iterate over our list of regions and use aws.StringValue to print the region name.
	instancesslice = instance
	wg.Add(len(regions.Regions))
	for _, region := range regions.Regions {
		go func(region *ec2.Region) {
			fmt.Println("Starting new go function with region: " + aws.StringValue(region.RegionName))
			queryInstances(aws.StringValue(region.RegionName))
		}(region)
	}

	wg.Wait()

	fmt.Println(instancesslice)
}

// GetInstances returns a list of Ec2instance structs that are currently running
func queryInstances(regionName string) {
	fmt.Println("Query Instance: " + regionName)
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
			mutex.Lock()
			*instancesslice = append(*instancesslice, isstruct)
			mutex.Unlock()
		}
	}
}
