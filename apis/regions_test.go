package apis

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/stretchr/testify/assert"
)

type describeRegionerMock struct {
	// this function field allows you to easily mock methods per test
	// https://medium.com/@matryer/meet-moq-easily-mock-interfaces-in-go-476444187d10
	DescribeRegionsFunc func(*ec2.DescribeRegionsInput) (*ec2.DescribeRegionsOutput, error)
}

func (mock describeRegionerMock) DescribeRegions(input *ec2.DescribeRegionsInput) (*ec2.DescribeRegionsOutput, error) {
	// call the mocked function field here
	return mock.DescribeRegionsFunc(input)
}

func Test_GetRegions(t *testing.T) {
	var regionName, endpoint string
	regionName = "regionName"
	endpoint = "endpoint"
	errorText := "someError"
	tests := map[string]struct {
		dr                   DescribeRegioner
		drInput              *ec2.DescribeRegionsInput
		drOutputResponse     *ec2.DescribeRegionsOutput
		errorResponse        error
		regionName, endpoint string
	}{
		"success": {
			dr: describeRegionerMock{
				// implement the function field
				DescribeRegionsFunc: func(input *ec2.DescribeRegionsInput) (*ec2.DescribeRegionsOutput, error) {
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
		"panic": {
			dr: describeRegionerMock{
				DescribeRegionsFunc: func(input *ec2.DescribeRegionsInput) (*ec2.DescribeRegionsOutput, error) {
					return &ec2.DescribeRegionsOutput{}, errors.New(errorText)
				},
			},
			drInput:          nil,
			drOutputResponse: &ec2.DescribeRegionsOutput{},
			errorResponse:    errors.New(errorText),
		},
	}

	for name, test := range tests {
		t.Logf("Running test case: %s", name)
		if name != "panic" {
			response := GetRegions(test.dr)
			assert.Equal(t, test.drOutputResponse, response)
		} else {
			assert.Panics(t, func() { GetRegions(test.dr) }, test.errorResponse.Error())
		}
	}
}
