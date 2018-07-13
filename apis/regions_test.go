package apis

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws/awserr"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/stretchr/testify/assert"
)

type describeRegionerMock struct {
	// this function field allows you to easily mock methods per test
	// https://medium.com/@matryer/meet-moq-easily-mock-interfaces-in-go-476444187d10
	DescribeRegionsFunc func(aws.Context, *ec2.DescribeRegionsInput, ...request.Option) (*ec2.DescribeRegionsOutput, error)
}

func (mock describeRegionerMock) DescribeRegionsWithContext(ctx aws.Context, input *ec2.DescribeRegionsInput, req ...request.Option) (*ec2.DescribeRegionsOutput, error) {
	// call the mocked function field here
	return mock.DescribeRegionsFunc(ctx, input, req...)
}

func Test_GetRegions(t *testing.T) {
	var regionName, endpoint string
	regionName = "regionName"
	endpoint = "endpoint"
	errorText := "someError"
	cancelText := "request cancelled"
	tests := map[string]struct {
		dr                   DescribeRegioner
		drInput              *ec2.DescribeRegionsInput
		ctx                  context.Context
		timeout              int
		drOutputResponse     *ec2.DescribeRegionsOutput
		errorResponse        error
		regionName, endpoint string
	}{
		"success": {
			dr: describeRegionerMock{
				// implement the function field
				DescribeRegionsFunc: func(ctx aws.Context, input *ec2.DescribeRegionsInput, req ...request.Option) (*ec2.DescribeRegionsOutput, error) {
					return &ec2.DescribeRegionsOutput{
						Regions: []*ec2.Region{
							{
								Endpoint:   &endpoint,
								RegionName: &regionName,
							},
						},
					}, nil
				},
			},
			drInput: nil,
			timeout: 100,
			ctx:     context.Background(),
			drOutputResponse: &ec2.DescribeRegionsOutput{
				Regions: []*ec2.Region{
					{
						Endpoint:   &endpoint,
						RegionName: &regionName,
					},
				},
			},
			errorResponse: nil,
		},
		"cancel": {
			dr: describeRegionerMock{
				DescribeRegionsFunc: func(ctx aws.Context, input *ec2.DescribeRegionsInput, req ...request.Option) (*ec2.DescribeRegionsOutput, error) {
					select {
					case <-ctx.Done():
						return nil, awserr.New(request.CanceledErrorCode, cancelText, nil)
					default:
						return &ec2.DescribeRegionsOutput{
							Regions: []*ec2.Region{
								{
									Endpoint:   &endpoint,
									RegionName: &regionName,
								},
							},
						}, nil
					}
				},
			},
			drInput:          nil,
			timeout:          0,
			ctx:              context.Background(),
			drOutputResponse: &ec2.DescribeRegionsOutput{},
			errorResponse:    errors.New(cancelText),
		},
		"panic": {
			dr: describeRegionerMock{
				DescribeRegionsFunc: func(ctx aws.Context, input *ec2.DescribeRegionsInput, req ...request.Option) (*ec2.DescribeRegionsOutput, error) {
					return &ec2.DescribeRegionsOutput{}, errors.New(errorText)
				},
			},
			drInput:          nil,
			timeout:          100,
			ctx:              context.Background(),
			drOutputResponse: &ec2.DescribeRegionsOutput{},
			errorResponse:    errors.New(errorText),
		},
	}

	for name, test := range tests {
		t.Logf("Running test case: %s", name)
		ctx, cancel := context.WithTimeout(test.ctx, 0*time.Second)
		defer cancel()
		if name == "cancel" || name == "panic" {
			assert.Panics(t, func() { GetRegions(ctx, test.dr) }, test.errorResponse.Error())
		} else {
			ctx, cancelFn := context.WithTimeout(test.ctx, time.Duration(test.timeout)*time.Second)
			defer cancelFn()
			response := GetRegions(ctx, test.dr)
			assert.Equal(t, test.drOutputResponse, response)
		}
	}
}
