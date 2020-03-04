package lib

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecr"
	log "github.com/sirupsen/logrus"
)

// ListNoncompliantECRRepositories lists all of the repositories within the current account and region
// and tracks repositories without image tag immutability enabled.
// TODO: What's the timeout bound on the number of requests made? Example: function timeout is 15 seconds
// but there is a limit of 10,000 repositories per region.
func (f *FunctionContainer) ListNoncompliantECRRepositories(ctx context.Context) ([]*ecr.Repository, error) {
	var repositories []*ecr.Repository
	input := &ecr.DescribeRepositoriesInput{}
	if err := input.Validate(); err != nil {
		return nil, err
	}
	pager := func(output *ecr.DescribeRepositoriesOutput, lastPage bool) bool {
		for _, r := range output.Repositories {
			if aws.StringValue(r.ImageTagMutability) == ecr.ImageTagMutabilityMutable {
				log.Infof("Repository %s does not have image tag immutability enabled", *r.RepositoryName)
				repositories = append(repositories, r)
			}
		}
		return lastPage
	}
	if err := f.ECR.DescribeRepositoriesPagesWithContext(ctx, input, pager); err != nil {
		return nil, err
	}
	if f.NotificationsEnabled {
		if err := f.PublishSNSMessage(ctx, repositories); err != nil {
			return nil, err
		}
	}
	return repositories, nil
}

// SetImageTagImmutability sets ImageTagImmutability on all noncompliant ECR repositories.
func (f *FunctionContainer) SetImageTagImmutability(ctx context.Context, repos []*ecr.Repository) error {
	for _, r := range repos {
		input := &ecr.PutImageTagMutabilityInput{
			ImageTagMutability: aws.String(ecr.ImageTagMutabilityImmutable),
			RepositoryName:     r.RepositoryName,
		}
		if err := input.Validate(); err != nil {
			return err
		}
		if _, err := f.ECR.PutImageTagMutabilityWithContext(ctx, input); err != nil {
			return err
		}
	}
	return nil
}
