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

invoke:
	aws lambda invoke --function-name ec2query outfile.txt

create_serverless_container:
	docker build -t $(BUILD_IMAGE) .

# use "make deploy AWS_ACCESS_KEY_ID=XXXX AWS_SECRET_ACCESS_KEY=XXX STAGE=dev REGION=us-east-1 SENDER=sender@example.com RECIPIENT=recipient@example.com"
# or if you've set the parameters as environment variables, simply use "make deploy"
# If you have the AWS shared credentials configured, those environment variables will be set for you.
deploy: build create_serverless_container
	eval $(go run $(PWD)/tools/awsconfig/main.go)
	docker run -it --rm -v $(PWD):/app -e AWS_ACCESS_KEY_ID=$(AWS_ACCESS_KEY_ID) -e AWS_SECRET_ACCESS_KEY=$(AWS_SECRET_ACCESS_KEY) $(BUILD_IMAGE) serverless deploy --stage=$(STAGE) --region=$(REGION) --sender=$(SENDER) --recipient=$(RECIPIENT) -v
undeploy: build create_serverless_container
	docker run -it --rm -v $(PWD):/app -e AWS_ACCESS_KEY_ID=$(AWS_ACCESS_KEY_ID) -e AWS_SECRET_ACCESS_KEY=$(AWS_SECRET_ACCESS_KEY) $(BUILD_IMAGE) serverless remove --stage=$(STAGE) --region=$(REGION) -v

clean:
	rm -f $(BINARY_NAME)
	rm -f $(ZIP_NAME)

deps:
	go get -v github.com/aws/aws-lambda-go/lambda
	go get -v github.com/aws/aws-sdk-go/aws
	go get -v github.com/aws/aws-sdk-go/aws/awserr
	go get -v github.com/aws/aws-sdk-go/aws/session
	go get -v github.com/aws/aws-sdk-go/service/ec2
	go get -v github.com/aws/aws-sdk-go/service/ses
	go get -v github.com/aws/aws-sdk-go/aws/credentials

test-deps:
	go get -v -t -d ./...

clean-testcache:
	go clean -testcache github.com/atssteve/ec2query/

test: test-deps	clean-testcache
	go test ./... -race -covermode=atomic

test-circleci: test-deps
	go test -race -covermode=atomic -coverprofile=coverage.txt ./...