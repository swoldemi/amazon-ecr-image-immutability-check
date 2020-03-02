// Package lib contains library units for the ecr-image-immutability-check Lambda function.
package lib

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/service/ecr/ecriface"
	log "github.com/sirupsen/logrus"
)

// Environment denotes different environments.
type Environment string

const (
	// Production denotes a production environment.
	Production Environment = "production"

	// Development denotes a development environment.
	Development Environment = "development"
)

// FunctionContainer contains the dependencies and business logic for the ecr-image-immutability-check Lambda function.
type FunctionContainer struct {
	Environment Environment
	ECR         ecriface.ECRAPI
}

// NewFunctionContainer creates a new FunctionContainer.
func NewFunctionContainer(ecrSvc ecriface.ECRAPI, env Environment) *FunctionContainer {
	log.Infof("Creating function container for environment: %v", env)
	return &FunctionContainer{
		Environment: env,
		ECR:         ecrSvc,
	}
}

// GetHandler returns the function handler for ecr-image-immutability-check.
func (f *FunctionContainer) GetHandler() func(context.Context, events.CloudWatchEvent) error {
	return func(ctx context.Context, event events.CloudWatchEvent) error {
		repos, err := f.ListNoncompliantECRRepositories(ctx)
		if err != nil {
			return err
		}
		if len(repos) == 0 {
			return nil
		}
		log.Infof("Found %d noncompliant ECR repositories\n", len(repos))
		if err := f.SetImageTagImmutability(ctx, repos); err != nil {
			return err
		}
		return nil
	}
}
