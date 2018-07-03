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
}

// APIMap stores the Region name are the key and a slice of Ec2Instance structs for use later
type APIMap map[string][]Ec2instance

var wg sync.WaitGroup
var mutex sync.Mutex
var instanceresult []APIMap

// GetInstances takes a pointer to a DescribeRegionsOutput to query all of the instances in a region and build on our Ec2instance struct
// The Ec2instance structs will be stored as a slice in a map to organized in fashion that can easily we looped over.GetInstances
//
// Example:
// map[us-east-1:[{i-030b2417f941cdbbf t2.micro 2018-06-28 17:11:20 +0000 UTC stopped aw-mac us-east-1}]]
func GetInstances(regions *ec2.DescribeRegionsOutput) []APIMap {
	// Iterate over our list of regions and use aws.StringValue to print the region name.
	fmt.Printf("Setting waitgroup count to: %d\n", len(regions.Regions)-1)
	wg.Add(len(regions.Regions))
	for _, region := range regions.Regions {
		go func(region *ec2.Region) {
			fmt.Println("Starting new go function with region: " + aws.StringValue(region.RegionName))
			m := queryInstances(aws.StringValue(region.RegionName))

			mutex.Lock()
			instanceresult = append(instanceresult, m)
			mutex.Unlock()
		}(region)
	}

	wg.Wait()

	return instanceresult
}

// GetInstances returns a list of Ec2instance structs that are currently running
func queryInstances(regionName string) APIMap {
	defer wg.Done()
	fmt.Println("Query Instance: " + regionName)
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(regionName),
	})
	if err != nil {
		panic(err)
	}

	ec2map := make(APIMap)
	var ec2slice []Ec2instance

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
			}
			ec2slice = append(ec2slice, isstruct)

		}
	}

	ec2map[regionName] = ec2slice

	return ec2map
}
