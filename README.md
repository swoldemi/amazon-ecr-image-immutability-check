![](https://codebuild.us-east-2.amazonaws.com/badges?uuid=eyJlbmNyeXB0ZWREYXRhIjoiZzVncHNUcDhJcVgrSE9aek9NS2J6Z3V3eW5BWkFOZTAxSTBWczVjSGNYVVE1WUpNcmpaM2RuZytXTnlSR1lNLzloeTJVeDJsMk9pQTY5anByd1hCakswPSIsIml2UGFyYW1ldGVyU3BlYyI6IlNXeVpSZ254ZnVkem0wTkgiLCJtYXRlcmlhbFNldFNlcmlhbCI6MX0%3D&branch=master)
[![][sar-logo]](https://serverlessrepo.aws.amazon.com/applications/arn:aws:serverlessrepo:us-east-1:273450712882:applications~amazon-ecr-image-immutability-check)


[sar-deploy]: https://img.shields.io/badge/Serverless%20Application%20Repository-Deploy%20Now-FF9900?logo=amazon%20aws&style=flat-square
[sar-logo]: https://img.shields.io/badge/Serverless%20Application%20Repository-View-FF9900?logo=amazon%20aws&style=flat-square

# Amazon ECR Image Immutability Check
>Enforce image tag immutability on all Elastic Container Registry repositories within an AWS account

The problem: As of March 2020, [AWS Config](https://aws.amazon.com/config/) does not support any custom or native integrations with ECR: https://docs.aws.amazon.com/config/latest/developerguide/resource-config-reference.html

The solution: Run a Serverless Application Repository app to automatically remediate and report on incompliant ECR repositories for you!

![architecture](https://raw.githubusercontent.com/swoldemi/amazon-ecr-image-immutability-check/master/screenshots/architecture.png)

## Usage
Prerequisites:
1. An AWS account in a region which supports ECR.
2. ECR repositories within your registry.

### Deploying the Lambda
It is recommended that you deploy this Lambda function directly from the AWS Serverless Application Repository. It is also possible to deploy this function using:
- The [SAM CLI](https://aws.amazon.com/serverless/sam/)
- CloudFormation via the [AWS CLI](https://aws.amazon.com/cli/)
- CloudFormation via the [CloudFormation management console](https://aws.amazon.com/cloudformation/)

This function has been made available in all 18 commercial AWS regions that support AWS SAR. It is also possible to deploy the Lambda function in the GovCloud and China regions, if you have access to those regions.

|Region                                        |Click and Deploy                                                                                                                                 |
|----------------------------------------------|-------------------------------------------------------------------------------------------------------------------------------------------------|
|**US East (Ohio) (us-east-2)**                |[![][sar-deploy]](https://deploy.serverlessrepo.app/us-east-2/?app=arn:aws:serverlessrepo:us-east-1:273450712882:applications/amazon-ecr-image-immutability-check)     |
|**US East (N. Virginia) (us-east-1)**         |[![][sar-deploy]](https://deploy.serverlessrepo.app/us-east-1/?app=arn:aws:serverlessrepo:us-east-1:273450712882:applications/amazon-ecr-image-immutability-check)     |
|**US West (N. California) (us-west-1)**       |[![][sar-deploy]](https://deploy.serverlessrepo.app/us-west-1/?app=arn:aws:serverlessrepo:us-east-1:273450712882:applications/amazon-ecr-image-immutability-check)     |
|**US West (Oregon) (us-west-2)**              |[![][sar-deploy]](https://deploy.serverlessrepo.app/us-west-2/?app=arn:aws:serverlessrepo:us-east-1:273450712882:applications/amazon-ecr-image-immutability-check)     |
|**Asia Pacific (Hong Kong) (ap-east-1)**      |[![][sar-deploy]](https://deploy.serverlessrepo.app/ap-east-1/?app=arn:aws:serverlessrepo:us-east-1:273450712882:applications/amazon-ecr-image-immutability-check)     |
|**Asia Pacific (Mumbai) (ap-south-1)**        |[![][sar-deploy]](https://deploy.serverlessrepo.app/ap-south-1/?app=arn:aws:serverlessrepo:us-east-1:273450712882:applications/amazon-ecr-image-immutability-check)    |
|**Asia Pacific (Seoul) (ap-northeast-2)**     |[![][sar-deploy]](https://deploy.serverlessrepo.app/ap-northeast-2/?app=arn:aws:serverlessrepo:us-east-1:273450712882:applications/amazon-ecr-image-immutability-check)|
|**Asia Pacific (Singapore)	(ap-southeast-1)** |[![][sar-deploy]](https://deploy.serverlessrepo.app/ap-southeast-1/?app=arn:aws:serverlessrepo:us-east-1:273450712882:applications/amazon-ecr-image-immutability-check)|
|**Asia Pacific (Sydney) (ap-southeast-2)**    |[![][sar-deploy]](https://deploy.serverlessrepo.app/ap-southeast-2/?app=arn:aws:serverlessrepo:us-east-1:273450712882:applications/amazon-ecr-image-immutability-check)|
|**Asia Pacific (Tokyo) (ap-northeast-1)**     |[![][sar-deploy]](https://deploy.serverlessrepo.app/ap-northeast-1?app=arn:aws:serverlessrepo:us-east-1:273450712882:applications/amazon-ecr-image-immutability-check) |
|**Canada (Central)	(ca-central-1)**           |[![][sar-deploy]](https://deploy.serverlessrepo.app/ca-central-1/?app=arn:aws:serverlessrepo:us-east-1:273450712882:applications/amazon-ecr-image-immutability-check)  |
|**EU (Frankfurt) (eu-central-1)**             |[![][sar-deploy]](https://deploy.serverlessrepo.app/eu-central-1/?app=arn:aws:serverlessrepo:us-east-1:273450712882:applications/amazon-ecr-image-immutability-check)  |
|**EU (Ireland)	(eu-west-1)**                  |[![][sar-deploy]](https://deploy.serverlessrepo.app/eu-west-1/?app=arn:aws:serverlessrepo:us-east-1:273450712882:applications/amazon-ecr-image-immutability-check)     |
|**EU (London) (eu-west-2)**                   |[![][sar-deploy]](https://deploy.serverlessrepo.app/eu-west-2/?app=arn:aws:serverlessrepo:us-east-1:273450712882:applications/amazon-ecr-image-immutability-check)     |
|**EU (Paris) (eu-west-3)**                    |[![][sar-deploy]](https://deploy.serverlessrepo.app/eu-west-3/?app=arn:aws:serverlessrepo:us-east-1:273450712882:applications/amazon-ecr-image-immutability-check)     |
|**EU (Stockholm) (eu-north-1)**               |[![][sar-deploy]](https://deploy.serverlessrepo.app/eu-north-1/?app=arn:aws:serverlessrepo:us-east-1:273450712882:applications/amazon-ecr-image-immutability-check)    |
|**Middle East (Bahrain) (me-south-1)**        |[![][sar-deploy]](https://deploy.serverlessrepo.app/me-south-1/?app=arn:aws:serverlessrepo:us-east-1:273450712882:applications/amazon-ecr-image-immutability-check)    |
|**South America (Sao Paulo) (sa-east-1)**     |[![][sar-deploy]](https://deploy.serverlessrepo.app/sa-east-1/?app=arn:aws:serverlessrepo:us-east-1:273450712882:applications/amazon-ecr-image-immutability-check)     |

### Configuration
1. Interval (required) - How often should the function run? Requires a valid Schedule Expression: https://docs.aws.amazon.com/lambda/latest/dg/tutorial-scheduled-events-schedule-expressions.html. Default is once a day (`rate(24 hours)`).
2. SNSTopicARN (optional) - The ARN of the Simple Notification Service topic to send incompliant finding messages to. No default.
3. AutoRemediate (optional) - Should this function automatically enable image tag immutability (ENABLED or DISABLED)? ENABLED by default.

If you would like to retrive notifications of repositories that have image immutability disabled create and subscribe to an SNS topic then pass the ARN of the topic in the SNSTopicARN parameter.

Example message contents (Email subscription):
```
The amazon-ecr-image-immutability-check Lambda function you deployed found some incompliant ECR repositories:
1. Repository: product-a/service-one
2. Repository: product-a/service-two
3. Repository: custom-internal-nginx
If auto-remediation is enabled, then these repositories will have Image Tag Immutability enabled and are now compliant until changed.
AWS Region: us-east-2
ECR Registry ID: 123456789012
Auto-Remediation: ENABLED
--
If you wish to stop receiving notifications from this topic, please click or visit the link below to unsubscribe:
https://sns.us-east-2.amazonaws.com/unsubscribe.html?SubscriptionArn=arn:aws:sns:us-east-2:123456789012:your-topic-name:random-topic-id&Endpoint=example@example.com

Please do not reply directly to this email. If you have any questions or comments regarding this email, please contact us at https://aws.amazon.com/support
```

### Test that it works
After your specified interval and interval unit (example: rate(5 minutes)), a CloudWatch event will trigger the Lambda function and scan your account for repositories that do not have image tag immutability enabled. If any are found, image tag immutability will be enabled.

## Contributing
Have an idea for a feature to enhance this serverless application? Open an [issue](https://github.com/swoldemi/amazon-ecr-image-immutability-check/issues) or [pull request](https://github.com/swoldemi/amazon-ecr-image-immutability-check/pulls)!

### Development
This application has been developed, built, and testing against [Go 1.13, Go 1.14](https://golang.org/dl/), the latest version of the [Serverless Application Model CLI](https://github.com/awslabs/aws-sam-cli), and the latest version of the [AWS CLI](https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-install.html). A [Makefile](./Makefile) has been provided for convenience.

```
make check
make test
make build
make sam-package
make sam-deploy
make sam-tail-logs
make destroy
```

## To Do
1. Integrate AWS Config when support for ECR repositories arrives.

## Note
It is also possible to prevent account users from changing the tag immutability setting by not granting them the `ecr:PutImageTagMutability` action.

## License
[MIT No Attribution (MIT-0)](https://spdx.org/licenses/MIT-0.html)
