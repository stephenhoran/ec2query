# specify any of the following parameters at the command line or set them as environment variables
# environment variables are used if command-line arguments are not provided
AWS_ACCESS_KEY_ID?=AWS_ACCESS_KEY_ID
AWS_SECRET_ACCESS_KEY?=AWS_SECRET_ACCESS_KEY
SENDER?=SENDER
RECIPIENT?=RECIPIENT
STAGE?=STAGE
REGION?=REGION
# end command-line parameters

GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=main
ZIP_NAME=$(BINARY_NAME).zip
BUILD_IMAGE=serverless-build:0.1

all: build push
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME) -v
push:
	export AWS_DEFAULT_REGION=us-east-2
	zip main.zip main
	aws lambda update-function-code --function-name ec2query --zip-file fileb://main.zip

create_serverless_container:
	docker build -t $(BUILD_IMAGE) .

# use "make deploy AWS_ACCESS_KEY_ID=XXXX AWS_SECRET_ACCESS_KEY=XXX STAGE=dev REGION=us-east-1 SENDER=sender@example.com RECIPIENT=recipient@example.com"
# or if you've set the parameters as environment variables, simply use "make deploy"
deploy: build create_serverless_container
	docker run -it --rm -v $(PWD):/app -e AWS_ACCESS_KEY_ID=$(AWS_ACCESS_KEY_ID) -e AWS_SECRET_ACCESS_KEY=$(AWS_SECRET_ACCESS_KEY) $(BUILD_IMAGE) serverless deploy --stage=$(STAGE) --region=$(REGION) --sender=$(SENDER) --recipient=$(RECIPIENT) -v
undeploy: build create_serverless_container
	docker run -it --rm -v $(PWD):/app -e AWS_ACCESS_KEY_ID=$(AWS_ACCESS_KEY_ID) -e AWS_SECRET_ACCESS_KEY=$(AWS_SECRET_ACCESS_KEY) $(BUILD_IMAGE) serverless remove --stage=$(STAGE) --region=$(REGION) -v

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


