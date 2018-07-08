package apis

import "github.com/aws/aws-sdk-go/service/ec2"

type DescribeRegioner interface {
	DescribeRegions(*ec2.DescribeRegionsInput) (*ec2.DescribeRegionsOutput, error)
}

// GetRegions will return a DescribeRegionsOutput stuct containing our
// slice of regions.
//
// It requires that you pass an already created session.
func GetRegions(sess DescribeRegioner) *ec2.DescribeRegionsOutput {
	result, err := sess.DescribeRegions(nil)
	if err != nil {
		panic(err)
	}

	return result
}
