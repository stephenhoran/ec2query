[![CircleCI](https://circleci.com/gh/atssteve/ec2query.svg?style=svg)](https://circleci.com/gh/atssteve/ec2query)
[![codecov](https://codecov.io/gh/atssteve/ec2query/branch/master/graph/badge.svg)](https://codecov.io/gh/atssteve/ec2query)
[![Go Report Card](https://goreportcard.com/badge/github.com/atssteve/ec2query)](https://goreportcard.com/report/github.com/atssteve/ec2query)

# EC2Query Lambda
This lambda will simply query all regions in AWS and email the recipient a list of what is running.

### Makefile
```
make deps # Will install all AWS dependancies
make build # Builds the Go binary in Linux to be compatible with Lambda
make push # Zips binary and updates current lambda
make clean # Cleans up directory
make # runs build and push
```
