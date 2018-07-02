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
}

func GetInstances(regions *ec2.DescribeRegionsOutput, htmlbody *string) {
	// Iterate over our list of regions and use aws.StringValue to print the region name.
	for _, region := range regions.Regions {
		var is []Ec2instance
		fmt.Println(aws.StringValue(region.RegionName))
		is = queryInstances(*region.RegionName)
		fmt.Println(len(is))
		if len(is) != 0 {
			*htmlbody = *htmlbody + "<h1>" + aws.StringValue(region.RegionName) + "</h1>"
			*htmlbody = *htmlbody + "<table border=\"1\"><th>Instance Names</th><th>Type</th><th>state</th><th>Launch Time</th><th>Key Name</th>"
			for _, i := range is {
				*htmlbody = *htmlbody + "<tr><td>" + i.Instanceid + "</td><td>" + i.Type + "</td><td>" + i.State + "</td><td>" + i.LaunchTime.Format("2006-01-02 15:04:05") + "</td><td>" + i.KeyName + "</td></tr>"
			}
			*htmlbody = *htmlbody + "</table>"
		}
	}
}

// GetInstances returns a list of Ec2instance structs that are currently running
func queryInstances(region string) []Ec2instance {
	var is []Ec2instance
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
			isstruct := Ec2instance{
				Instanceid: aws.StringValue(instances.InstanceId),
				Type:       aws.StringValue(instances.InstanceType),
				LaunchTime: instances.LaunchTime,
				State:      aws.StringValue(instances.State.Name),
				KeyName:    aws.StringValue(instances.KeyName),
			}
			is = append(is, isstruct)
		}
	}

	return is
}
