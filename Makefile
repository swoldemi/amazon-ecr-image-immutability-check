all : tmpl check build sam-package sam-deploy sam-tail-logs
.PHONY: all

S3_BUCKET ?= swoldemi-tmp
DEFAULT_STACK_NAME ?= lambda-amazon-ecr-image-immutability-check
DEFAULT_REGION ?= us-east-2

GOBIN := $(GOPATH)/bin
GOIMPORTS := $(GOBIN)/goimports
GOLANGCILINT := $(GOBIN)/golangci-lint
GOREPORTCARDCLI := $(GOBIN)/goreportcard-cli
GOMETALINTER := $(GOBIN)/gometalinter

.PHONY: build
build: clean
	 go build -v -a -installsuffix cgo -tags netgo -ldflags '-w -extldflags "-static"' main.go

.PHONY: test
test: clean
	go test -v -race -timeout 30s -count=1 -coverprofile=profile.out ./...

# Static code analysis tooling and checks
.PHONY: check
check:
	goimports -w -l -e .
	gofmt -s -w .
	golangci-lint run ./... \
		-E goconst \
		-E gocyclo \
		-E gosec  \
		-E gofmt \
		-E maligned \
		-E misspell \
		-E nakedret \
		-E unconvert \
		-E unparam \
		-E dupl
	goreportcard-cli -v -t 90

.PHONY: tmpl
tmpl: 
	cfn-lint template.yaml

.PHONY: sam-package
sam-package:
	sam package --template-file template.yaml --s3-bucket $(S3_BUCKET) --output-template-file packaged.yaml

.PHONY: sam-deploy
sam-deploy:
	sam deploy \
	--region $(DEFAULT_REGION) \
	--template-file ./packaged.yaml \
	--parameter-overrides Interval="rate(1\ minute)" SNSTopicARN=arn:aws:sns:us-east-2:273450712882:codestar-notifications-distris \
	--stack-name $(DEFAULT_STACK_NAME) \
	--capabilities CAPABILITY_IAM
	aws --region $(DEFAULT_REGION) cloudformation describe-stacks --stack-name $(DEFAULT_STACK_NAME) --query 'Stacks[0].Outputs'

.PHONY: sam-publish
sam-publish:
	sam publish --region us-east-1 --template packaged.yaml

.PHONY: stack-describe 
stack-describe:
	aws --region $(DEFAULT_REGION) cloudformation describe-stacks --stack-name $(DEFAULT_STACK_NAME) --query 'Stacks[0].Outputs'

.PHONY: sam-tail-logs
sam-tail-logs:
	sam logs --name amazon-ecr-image-immutability-check --tail

.PHONY: destroy
destroy: clean
	aws --region $(DEFAULT_REGION) cloudformation delete-stack --stack-name $(DEFAULT_STACK_NAME)

.PHONY: update
update:
	go get $(shell go list -f "{{if not (or .Main .Indirect)}}{{.Path}}{{end}}" -m all)
	go mod tidy

.PHONY: clean
clean:
	rm -f main packaged.yaml profile.out


.PHONY: sar-public
sar-public:
	# Use this to make your SAR application public to all AWS accounts
	aws serverlessrepo put-application-policy \
		--region us-east-1 \
		--application-id arn:aws:serverlessrepo:us-east-1:273450712882:applications/amazon-ecr-image-immutability-check \
		--statements Principals=*,Actions=Deploy
