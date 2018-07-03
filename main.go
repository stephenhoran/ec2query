package main

import (
	"os"
	"runtime"

	"github.com/atssteve/ec2query/apis"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

const (
	// CharSet for SES email
	CharSet = "UTF-8"
)

// Handle used by the lambda
func Handle() {
	var instances []apis.Ec2instance
	runtime.GOMAXPROCS(runtime.NumCPU())

	sess, err := session.NewSession(&aws.Config{Region: aws.String(os.Getenv("REGION"))})
	if err != nil {
		panic(err)
	}

	// Starting a new EC2 client.
	ec2session := ec2.New(sess)

	// Describe all of the regions in AWS. Returns a type *DescribeRegionsOutput.
	resultRegions := apis.GetRegions(ec2session)

	// Query all regions for all instance types and append to email
	apis.GetInstances(resultRegions, &instances)

	//fmt.Println(artifacts)

	// 	date := time.Now()
	// 	svc := ses.New(sess)
	// 	// Assemble the email.
	// 	input := &ses.SendEmailInput{
	// 		Destination: &ses.Destination{
	// 			CcAddresses: []*string{},
	// 			ToAddresses: []*string{
	// 				aws.String(os.Getenv("RECIPIENT")),
	// 			},
	// 		},
	// 		Message: &ses.Message{
	// 			Body: &ses.Body{
	// 				Html: &ses.Content{
	// 					Charset: aws.String(CharSet),
	// 					Data:    aws.StringValue(HTMLBody),
	// 				},
	// 			},
	// 			Subject: &ses.Content{
	// 				Charset: aws.String(CharSet),
	// 				Data:    aws.String("AWS Report " + date.Format("01-02")),
	// 			},
	// 		},
	// 		Source: aws.String(os.Getenv("SENDER")),
	// 		// Uncomment to use a configuration set
	// 		//ConfigurationSetName: aws.String(ConfigurationSet),
	// 	}

	// 	// Attempt to send the email.
	// 	result, err := svc.SendEmail(input)
	// 	if err != nil {
	// 		if aerr, ok := err.(awserr.Error); ok {
	// 			switch aerr.Code() {
	// 			case ses.ErrCodeMessageRejected:
	// 				fmt.Println(ses.ErrCodeMessageRejected, aerr.Error())
	// 			case ses.ErrCodeMailFromDomainNotVerifiedException:
	// 				fmt.Println(ses.ErrCodeMailFromDomainNotVerifiedException, aerr.Error())
	// 			case ses.ErrCodeConfigurationSetDoesNotExistException:
	// 				fmt.Println(ses.ErrCodeConfigurationSetDoesNotExistException, aerr.Error())
	// 			default:
	// 				fmt.Println(aerr.Error())
	// 			}
	// 		} else {
	// 			// Print the error, cast err to awserr.Error to get the Code and
	// 			// Message from an error.
	// 			fmt.Println(err.Error())
	// 		}

	// 		return
}

// 	fmt.Println(result)
// }

func main() {
	Handle()
}
