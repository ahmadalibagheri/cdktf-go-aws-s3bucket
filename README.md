# cdktf-go-aws-s3bucket

The Cloud Development Kit for Terraform (CDKTF) allows you to define your infrastructure in a familiar programming language such as TypeScript, Python, Go, C#, or Java.

In this tutorial, you will provision an EC2 instance on AWS using your preferred programming language.

## Prerequisites

* [Terraform](https://www.terraform.io/downloads) >= v1.0
* [CDK for Terraform](https://learn.hashicorp.com/tutorials/terraform/cdktf-install) >= v0.8
* A [Terraform Cloud](https://app.terraform.io/) account, with [CLI authentication](https://learn.hashicorp.com/tutorials/terraform/cloud-login) configured
* [an AWS account](https://portal.aws.amazon.com/billing/signup?nc2=h_ct&src=default&redirect_url=https%3A%2F%2Faws.amazon.com%2Fregistration-confirmation#/start)
* AWS Credentials [configured for use with Terraform](https://registry.terraform.io/providers/hashicorp/aws/latest/docs#authentication)


Credentials can be provided by using the AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY, and optionally AWS_SESSION_TOKEN environment variables. The region can be set using the AWS_REGION or AWS_DEFAULT_REGION environment variables.

```shell
$ export AWS_ACCESS_KEY_ID="anaccesskey"
$ export AWS_SECRET_ACCESS_KEY="asecretkey"
$ export AWS_REGION="us-west-2"
```

## Install project dependencies

```shell
mkdir cdktf-go-aws-s3bucket
cd cdktf-go-aws-s3bucket
cdktf init --template="go"
```

## Install AWS provider
Open `cdktf.json` in your text editor, and add `aws` as one of the Terraform providers that you will use in the application.
```JSON
{
  "language": "go",
  "app": "go run main.go",
  "codeMakerOutput": "generated",
  "projectId": "02f2d864-a2f2-49e8-ab52-b472e233755e",
  "sendCrashReports": "false",
  "terraformProviders": [
	 "hashicorp/aws@~> 3.67.0"
  ],
  "terraformModules": [],
  "context": {
    "excludeStackIdFromLogicalIds": "true",
    "allowSepCharsInLogicalIds": "true"
  }
}
```
Run `cdktf get` to install the AWS provider you added to `cdktf.json`.
```SHELL
cdktf get
```

CDKTF uses a library called `jsii` to allow Go code to interact with CDK, 
which is written in TypeScript. 
Ensure that the jsii runtime is installed by running `go mod tidy`.

```SHELL
go mod tidy
```

## Define your CDK for Terraform Application

Replace the contents of main.py with the following code for a new Golang application

```golang
package main

import (
	"fmt"

	"cdk.tf/go/stack/generated/hashicorp/aws"
	"cdk.tf/go/stack/generated/hashicorp/aws/s3"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/hashicorp/terraform-cdk-go/cdktf"
)

func NewMyStack(scope constructs.Construct, id string) cdktf.TerraformStack {
	stack := cdktf.NewTerraformStack(scope, &id)
	aws.NewAwsProvider(stack, jsii.String("AWS"), &aws.AwsProviderConfig{
		Region: jsii.String("us-east-1"),
	})

	bucketName := "cdktf-golang-demo-us-east-1"
	lifeCycle := &[]s3.S3BucketLifecycleRule{
		{
			Enabled:                            true,
			Id:                                 jsii.String("abort-multipart"),
			Prefix:                             jsii.String("/"),
			AbortIncompleteMultipartUploadDays: jsii.Number(7),
		},
		{
			Enabled: true,
			Transition: &[]s3.S3BucketLifecycleRuleTransition{
				{
					StorageClass: jsii.String("STANDARD_IA"),
					Days:         jsii.Number(30),
				},
			},
		},
		{
			Enabled: true,
			Transition: &[]s3.S3BucketLifecycleRuleNoncurrentVersionTransition{
				{
					StorageClass: jsii.String("STANDARD_IA"),
					Days:         jsii.Number(30),
				},
			},
		},
		{
			Enabled: false,
			Transition: &[]s3.S3BucketLifecycleRuleTransition{
				{
					StorageClass: jsii.String("ONEZONE_IA"),
					Days:         jsii.Number(90),
				},
			},
		},
		{
			Enabled: false,
			Transition: &[]s3.S3BucketLifecycleRuleNoncurrentVersionTransition{
				{
					StorageClass: jsii.String("ONEZONE_IA"),
					Days:         jsii.Number(90),
				},
			},
		},
		{
			Enabled: false,
			Transition: &[]s3.S3BucketLifecycleRuleTransition{
				{
					StorageClass: jsii.String("GLACIER"),
					Days:         jsii.Number(365),
				},
			},
		},
		{
			Enabled: false,
			Transition: &[]s3.S3BucketLifecycleRuleNoncurrentVersionTransition{
				{
					StorageClass: jsii.String("GLACIER"),
					Days:         jsii.Number(365),
				},
			},
		},
	}

	// define tags
	tags := &map[string]*string{
		"Team":    jsii.String("Devops"),
		"Company": jsii.String("Your compnay"),
	}

	policy := fmt.Sprintf(`{
        "Version": "2012-10-17",
        "Statement": [
          {
            "Sid": "PublicReadGetObject",
            "Effect": "Allow",
            "Principal": "*",
            "Action": [
              "s3:GetObject"
            ],
            "Resource": [
              "arn:aws:s3:::%s/*"
            ]
          }
        ]
      }`, bucketName)

	bucket := s3.NewS3Bucket(stack, jsii.String("bucket"), &s3.S3BucketConfig{
		Bucket:        jsii.String(bucketName),
		LifecycleRule: lifeCycle,
		Tags:          tags,
		Policy:        jsii.String(policy),
	})

	// define outputs
	cdktf.NewTerraformOutput(stack, jsii.String("s3-id"), &cdktf.TerraformOutputConfig{
		Value: bucket.Id(),
	})

	cdktf.NewTerraformOutput(stack, jsii.String("s3-arn"), &cdktf.TerraformOutputConfig{
		Value: bucket.Arn(),
	})

	return stack
}

func main() {
	app := cdktf.NewApp(nil)

	stack := NewMyStack(app, "cdktf-go-aws-s3bucket")
	cdktf.NewRemoteBackend(stack, &cdktf.RemoteBackendProps{
		Hostname:     jsii.String("app.terraform.io"),
		Organization: jsii.String("jigsaw373"),
		Workspaces:   cdktf.NewNamedRemoteWorkspace(jsii.String("cdktf-go-aws-s3bucket")),
	})

	app.Synth()
}

```
## Provision infrastructure
```shell
cdktf deploy
```
After the instance is created, visit the AWS EC2 Dashboard.

## Clean up your infrastructure
```shell
cdktf destroy
```
