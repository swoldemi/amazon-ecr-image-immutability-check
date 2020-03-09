// Package lib contains library units for the amazon-ecr-image-immutability-check Lambda function.
package lib

import (
	"context"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/service/ecr/ecriface"
	"github.com/aws/aws-sdk-go/service/sns/snsiface"
	log "github.com/sirupsen/logrus"
)

// Environment denotes different environments.
type Environment string

// AutoRemediationStatus the status of autoremediation.
type AutoRemediationStatus string

const (
	// Production denotes a production environment.
	Production Environment = "production"

	// Development denotes a development environment.
	Development Environment = "development"

	// Enabled denoted that auto-remediation is enabled.
	Enabled AutoRemediationStatus = "ENABLED"

	// Disabled denoted that auto-remediation is disabled.
	Disabled AutoRemediationStatus = "DISABLED"

)

// FunctionContainer contains the dependencies and business logic for the amazon-ecr-image-immutability-check Lambda function.
type FunctionContainer struct {
	Environment            Environment
	ECR                    ecriface.ECRAPI
	SNS                    snsiface.SNSAPI
	TopicARN               string
	NotificationsEnabled   bool
	AutoRemediationStatus AutoRemediationStatus
	AutoRemediationEnabled bool
}

// NewFunctionContainer creates a new FunctionContainer.
func NewFunctionContainer(ecrSvc ecriface.ECRAPI, snsSvc snsiface.SNSAPI, env Environment) *FunctionContainer {
	log.Infof("Creating function container for environment: %v", env)
	return &FunctionContainer{
		Environment: env,
		ECR:         ecrSvc,
		SNS:         snsSvc,
	}
}

// GetHandler returns the function handler for amazon-ecr-image-immutability-check.
func (f *FunctionContainer) GetHandler() func(context.Context, events.CloudWatchEvent) error {
	topicARN := os.Getenv("SNS_TOPIC_ARN")
	if topicARN != "" {
		f.NotificationsEnabled = true
		f.TopicARN = topicARN
	}
	if os.Getenv("AUTO_REMEDIATE") == "ENABLED" {
		f.AutoRemediationEnabled = true
		f.AutoRemediationStatus = Enabled
	}
	if !f.AutoRemediationEnabled {
		f.AutoRemediationStatus = Disabled
	}
	return func(ctx context.Context, event events.CloudWatchEvent) error {
		repos, err := f.ListIncompliantECRRepositories(ctx)
		if err != nil {
			return err
		}
		if len(repos) == 0 {
			return nil
		}
		log.Infof("Found %d incompliant ECR repositories\n", len(repos))
		if err := f.SetImageTagImmutability(ctx, repos); err != nil {
			return err
		}
		return nil
	}
}
