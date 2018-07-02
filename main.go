package main

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/atssteve/ec2query/apis"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ses"
)

const (
	// CharSet for SES email
	CharSet = "UTF-8"
)

// Main used by the lambda
func main() {
	var HTMLBody string
	runtime.GOMAXPROCS(runtime.NumCPU())

	sess, err := session.NewSession(&aws.Config{Region: aws.String(os.Getenv("REGION"))})
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

	// Iterate over our list of regions and use aws.StringValue to print the region name.
	for _, region := range resultRegions.Regions {
		var is []apis.Ec2instance
		fmt.Println(aws.StringValue(region.RegionName))
		is = apis.GetInstances(*region.RegionName)
		fmt.Println(len(is))
		if len(is) != 0 {
			HTMLBody = HTMLBody + "<h1>" + aws.StringValue(region.RegionName) + "</h1>"
			HTMLBody = HTMLBody + "<table border=\"1\"><th>Instance Names</th><th>Type</th><th>state</th><th>Launch Time</th><th>Key Name</th>"
			for _, i := range is {
				HTMLBody = HTMLBody + "<tr><td>" + i.Instanceid + "</td><td>" + i.Type + "</td><td>" + i.State + "</td><td>" + i.LaunchTime.Format("2006-01-02 15:04:05") + "</td><td>" + i.KeyName + "</td></tr>"
			}
			HTMLBody = HTMLBody + "</table>"
		}
		//artifacts[aws.StringValue(region.RegionName)] = is
	}
	//fmt.Println(artifacts)

	date := time.Now()
	svc := ses.New(sess)
	// Assemble the email.
	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses: []*string{},
			ToAddresses: []*string{
				aws.String(os.Getenv("RECIPIENT")),
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
				Data:    aws.String("AWS Report " + date.Format("01-02")),
			},
		},
		Source: aws.String(os.Getenv("SENDER")),
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
