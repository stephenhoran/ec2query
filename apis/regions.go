package apis

import "github.com/aws/aws-sdk-go/service/ec2"

// GetRegions will return a DescribeRegionsOutput stuct containing our
// slice of regions.
//
// It requires that you pass an already created session.
func GetRegions(sess *ec2.EC2) *ec2.DescribeRegionsOutput {
	result, err := sess.DescribeRegions(nil)
	if err != nil {
		panic(err)
	}

	return result
}
