package lib

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/aws/aws-sdk-go/service/sns"
	log "github.com/sirupsen/logrus"
)

const (
	header = `The ecr-image-immutability-check Lambda function you deployed found some incompliant ECR repositories:`

	footer = `
These repositories have had Image Tag Immutability enabled and are now compliant until changed.
AWS Region: %s
ECR Registry ID: %s
`
)

// ConstructMessage constructurs the message sent to the SNS topic using
// the incompliant repositories that were found.
func ConstructMessage(repos []*ecr.Repository) (string, error) {
	message := new(strings.Builder)
	if _, err := message.WriteString(header); err != nil {
		return "", err
	}
	for i, r := range repos {
		part := fmt.Sprintf("\n%d. Repository: %s", i+1, *r.RepositoryName)
		if _, err := message.WriteString(part); err != nil {
			return "", err
		}
	}
	message.WriteString(fmt.Sprintf(footer, os.Getenv("AWS_REGION"), *repos[0].RegistryId))
	return message.String(), nil
}

// PublishSNSMessage reports incompliant repository findings to an SNS topic.
func (f *FunctionContainer) PublishSNSMessage(ctx context.Context, repos []*ecr.Repository) error {
	message, err := ConstructMessage(repos)
	if err != nil {
		return err
	}
	log.Infof("Sending message: %v", message)
	input := &sns.PublishInput{
		Message:  aws.String(message),
		TopicArn: aws.String(f.TopicARN),
	}
	if err := input.Validate(); err != nil {
		return err
	}
	output, err := f.SNS.PublishWithContext(ctx, input)
	if err != nil {
		return err
	}
	log.Infof("Got message ID: %s", *output.MessageId)
	return nil
}
