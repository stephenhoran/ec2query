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
