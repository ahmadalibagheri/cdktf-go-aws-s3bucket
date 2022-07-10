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
