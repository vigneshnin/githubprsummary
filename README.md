# sailpoint-githubprsummary
An AWS Serverless solution that sends Github PR summary as email.
This creates a Lambda function triggered every 7 days using EventBridge rule. It uses SNS to send email notifications.

## Pre-requisites
* AWS SAM CLI
* AWS Account
* AWS Profile with Access and Secret Key or SSO
* Go 1.20.x or higher

## Making it work
The following steps can be followed to make it work.

### AWS Side
The samconfig.toml file contains the configurations required for the SAM CLI. We are interested in the prod environment which has already been configured.
The configurations like region, profile, s3 bucket etc. should be changed as per the requirement.

### Application side
The application only needs the repo owner and the repo name. This can be updated in the samconfig.toml file in parameter_overrides section

### Deployment
Run the following command to build the application
```
sam build
```

Run the following command to deploy the application
```
sam deploy --config-env prod
```

### Configuring the Emails
Once the application is deployed, an SNS will be created. The emails which require the notification can be added as the SNS subscription.
