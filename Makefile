GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=main
ZIP_NAME=$(BINARY_NAME).zip

all: build push
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME) -v
push:
	export AWS_DEFAULT_REGION=us-east-2
	zip main.zip main
	aws lambda update-function-code --function-name ec2query --zip-file fileb://main.zip
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


