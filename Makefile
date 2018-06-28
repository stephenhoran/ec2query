GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=main
ZIP_NAME=$(BINARY_NAME).zip
BUILD_IMAGE=serverless-build:0.1
# read AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY from command line
# use environment variables if command-line arguments are not provided
AWS_ACCESS_KEY_ID?=AWS_ACCESS_KEY_ID
AWS_SECRET_ACCESS_KEY?=AWS_SECRET_ACCESS_KEY

all: build push
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME) -v
push:
	aws lambda update-function-code --function-name ec2query --zip-file fileb://main.zip

# for deploy-dev/undeploy-dev, use: make push AWS_ACCESS_KEY_ID=XXXX AWS_SECRET_ACCESS_KEY=XXX
# - OR -
# set the AWS_ACCESS_KEY_ID=XXXX AWS_SECRET_ACCESS_KEY=XXX environment variables locally and just use "make push"
create_serverless_container:
	docker build -t $(BUILD_IMAGE) .
deploy-dev: build create_serverless_container
	docker run -it --rm -v $(PWD):/app -e AWS_ACCESS_KEY_ID=$(AWS_ACCESS_KEY_ID) -e AWS_SECRET_ACCESS_KEY=$(AWS_SECRET_ACCESS_KEY) $(BUILD_IMAGE) serverless deploy --stage=dev --region=us-east-1 -v
undeploy-dev: build create_serverless_container
	docker run -it --rm -v $(PWD):/app -e AWS_ACCESS_KEY_ID=$(AWS_ACCESS_KEY_ID) -e AWS_SECRET_ACCESS_KEY=$(AWS_SECRET_ACCESS_KEY) $(BUILD_IMAGE) serverless remove --stage=dev --region=us-east-1 -v

clean:
	rm -f $(BINARY_NAME)
	rm -f $(ZIP_NAME)
deps:
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ses"


