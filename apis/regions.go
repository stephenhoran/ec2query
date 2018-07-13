package apis

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type DescribeRegioner interface {
	DescribeRegionsWithContext(aws.Context, *ec2.DescribeRegionsInput, ...request.Option) (*ec2.DescribeRegionsOutput, error)
}

// GetRegions will return a DescribeRegionsOutput stuct containing our
// slice of regions.
//
// It requires that you pass an already created session.
func GetRegions(ctx context.Context, sess DescribeRegioner) *ec2.DescribeRegionsOutput {

	result, err := sess.DescribeRegionsWithContext(ctx, nil)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok && aerr.Code() == request.CanceledErrorCode {
			fmt.Println("request's context canceled,", err)
		}
		panic(err)
	}

	return result

}
