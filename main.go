package main

import (
	"fmt"
	"runtime"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ses"
)

const (
	Sender    = "steve.horan@theatsgroup.com"
	Recipient = "steve.horan@theatsgroup.com"
	Subject   = "AWS Report"
	CharSet   = "UTF-8"
)

func main() {
	var HTMLBody string
	runtime.GOMAXPROCS(runtime.NumCPU())

	sess, err := session.NewSession(&aws.Config{Region: aws.String("us-east-1")})
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

	//artifacts := make(map[string][]string)

	// Grab the length of the slice of regions and create a WaitGroup for this.
	// Iterate over our list of regions and use aws.StringValue to print the region name.
	for _, region := range resultRegions.Regions {
		var is []string
		fmt.Println(aws.StringValue(region.RegionName))
		is = getInstances(*region.RegionName)
		fmt.Println(len(is))
		if len(is) != 0 {
			HTMLBody = HTMLBody + "<h1>" + aws.StringValue(region.RegionName) + "</h1>"
			HTMLBody = HTMLBody + "<table border=\"1\"><th>Instance Names</th>"
			for _, i := range is {
				HTMLBody = HTMLBody + "<tr><td>" + i + "</td></tr>"
			}
			HTMLBody = HTMLBody + "</table>"
		}
		//artifacts[aws.StringValue(region.RegionName)] = is
	}
	//fmt.Println(artifacts)

	svc := ses.New(sess)
	// Assemble the email.
	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses: []*string{},
			ToAddresses: []*string{
				aws.String(Recipient),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String(CharSet),
					Data:    aws.String(HTMLBody),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String(CharSet),
				Data:    aws.String(Subject),
			},
		},
		Source: aws.String(Sender),
		// Uncomment to use a configuration set
		//ConfigurationSetName: aws.String(ConfigurationSet),
	}

	// Attempt to send the email.
	result, err := svc.SendEmail(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ses.ErrCodeMessageRejected:
				fmt.Println(ses.ErrCodeMessageRejected, aerr.Error())
			case ses.ErrCodeMailFromDomainNotVerifiedException:
				fmt.Println(ses.ErrCodeMailFromDomainNotVerifiedException, aerr.Error())
			case ses.ErrCodeConfigurationSetDoesNotExistException:
				fmt.Println(ses.ErrCodeConfigurationSetDoesNotExistException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}

		return
	}

	fmt.Println(result)
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
