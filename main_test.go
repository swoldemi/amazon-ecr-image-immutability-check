package main

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/aws/aws-sdk-go/service/ecr/ecriface"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
	"github.com/swoldemi/ecr-image-immutability-check/pkg/lib"
)

type mockECRClient struct {
	mock.Mock
	ecriface.ECRAPI
}

// DescribeRepositoriesPagesWithContext mocks the DescribeRepositoriesPagesWithContext ECR API endpoint.
func (_m *mockECRClient) DescribeRepositoriesPagesWithContext(ctx aws.Context, input *ecr.DescribeRepositoriesInput, fn func(*ecr.DescribeRepositoriesOutput, bool) bool, opts ...request.Option) error {
	log.Debug("Mocking DescribeRepositoriesPagesWithContext API")
	args := _m.Called(ctx, input, fn)
	return args.Error(0)
}

// PutImageTagMutabilityWithContext mocks the PutImageTagMutabilityWithContext ECR API endpoint.
func (_m *mockECRClient) PutImageTagMutabilityWithContext(ctx aws.Context, input *ecr.PutImageTagMutabilityInput, opts ...request.Option) (*ecr.PutImageTagMutabilityOutput, error) {
	log.Debug("Mocking PutImageTagMutabilityWithContext API")
	args := _m.Called(ctx, input)
	return args.Get(0).(*ecr.PutImageTagMutabilityOutput), args.Error(1)
}

var defaultEvent = events.CloudWatchEvent{
	Version:    "0",
	ID:         "89d1a02d-5ec7-412e-82f5-13505f849b41",
	DetailType: "Scheduled Event",
	Source:     "aws.events",
	AccountID:  "123456789012",
	Time:       time.Now(),
	Region:     "us-east-1",
	Resources:  []string{"arn:aws:events:us-east-1:123456789012:rule/SampleRule"},
	Detail:     json.RawMessage{},
}

func TestHandler(t *testing.T) {
	testRepos := []*ecr.Repository{
		{RepositoryName: aws.String("test-repo-one"), ImageTagMutability: aws.String(ecr.ImageTagMutabilityMutable)},
		{RepositoryName: aws.String("test-repo-two"), ImageTagMutability: aws.String(ecr.ImageTagMutabilityMutable)},
	}
	ecrSvc := new(mockECRClient)
	ecrSvc.On("DescribeRepositoriesPagesWithContext",
		context.Background(),
		&ecr.DescribeRepositoriesInput{},
		mock.AnythingOfType("func(*ecr.DescribeRepositoriesOutput, bool) bool"),
	).Return(nil).Run(func(args mock.Arguments) {
		arg := args.Get(2).(func(*ecr.DescribeRepositoriesOutput, bool) bool)
		arg(&ecr.DescribeRepositoriesOutput{Repositories: testRepos}, true)
	})

	for _, r := range testRepos {
		input := &ecr.PutImageTagMutabilityInput{
			ImageTagMutability: aws.String(ecr.ImageTagMutabilityImmutable),
			RepositoryName:     r.RepositoryName,
		}
		ecrSvc.On("PutImageTagMutabilityWithContext", context.Background(), input).
			Return(&ecr.PutImageTagMutabilityOutput{}, nil)
	}

	f := lib.NewFunctionContainer(ecrSvc, lib.Development)
	if err := f.GetHandler()(context.Background(), defaultEvent); err != nil {
		t.Fatalf("Error invoking handler: %v\n", err)
	}
	ecrSvc.AssertExpectations(t)
}
