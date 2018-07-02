package apis

import (
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

// GetInstances returns a list of Ec2instance structs that are currently running
func GetInstances(region string) []Ec2instance {
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
