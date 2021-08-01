package cmd

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type S3CreateBucketAPI interface {
	CreateBucket(ctx context.Context,
		params *s3.CreateBucketInput,
		optFns ...func(*s3.Options)) (*s3.CreateBucketOutput, error)
}

func createS3Bucket(c context.Context, api S3CreateBucketAPI, bucketName string, region string) (*s3.CreateBucketOutput, error) {
	in := &s3.CreateBucketInput{
		Bucket: &bucketName,
		CreateBucketConfiguration: &types.CreateBucketConfiguration{
			LocationConstraint: types.BucketLocationConstraint(region),
		},
	}
	return api.CreateBucket(c, in)
}

type S3PutPublicAccessBlockAPI interface {
	PutPublicAccessBlock(ctx context.Context,
		params *s3.PutPublicAccessBlockInput,
		optFns ...func(*s3.Options)) (*s3.PutPublicAccessBlockOutput, error)
}

func enableAllPublicAccessBlock(c context.Context, api S3PutPublicAccessBlockAPI, bucketName string) (*s3.PutPublicAccessBlockOutput, error) {
	in := &s3.PutPublicAccessBlockInput{
		Bucket: &bucketName,
		PublicAccessBlockConfiguration: &types.PublicAccessBlockConfiguration{
			BlockPublicAcls:       true,
			BlockPublicPolicy:     true,
			IgnorePublicAcls:      true,
			RestrictPublicBuckets: true,
		},
	}
	return api.PutPublicAccessBlock(c, in)
}

type S3PutBucketEncryptionAPI interface {
	PutBucketEncryption(ctx context.Context,
		params *s3.PutBucketEncryptionInput,
		optFns ...func(*s3.Options)) (*s3.PutBucketEncryptionOutput, error)
}

func enableBucketEncryptionAES256(c context.Context, api S3PutBucketEncryptionAPI, bucketName string) (*s3.PutBucketEncryptionOutput, error) {
	in := &s3.PutBucketEncryptionInput{
		Bucket: &bucketName,
		ServerSideEncryptionConfiguration: &types.ServerSideEncryptionConfiguration{
			Rules: []types.ServerSideEncryptionRule{
				{
					ApplyServerSideEncryptionByDefault: &types.ServerSideEncryptionByDefault{
						SSEAlgorithm: types.ServerSideEncryptionAes256,
					},
				},
			},
		},
	}
	return api.PutBucketEncryption(c, in)
}

type S3PutBucketVersioningAPI interface {
	PutBucketVersioning(ctx context.Context,
		params *s3.PutBucketVersioningInput,
		optFns ...func(*s3.Options)) (*s3.PutBucketVersioningOutput, error)
}

func enableBucketVersioning(c context.Context, api S3PutBucketVersioningAPI, bucketName string) (*s3.PutBucketVersioningOutput, error) {
	in := &s3.PutBucketVersioningInput{
		Bucket: &bucketName,
		VersioningConfiguration: &types.VersioningConfiguration{
			Status: types.BucketVersioningStatusEnabled,
		},
	}
	return api.PutBucketVersioning(c, in)
}

type S3GetbucketLocation interface {
	GetBucketLocation(ctx context.Context,
		params *s3.GetBucketLocationInput,
		optFns ...func(*s3.Options)) (*s3.GetBucketLocationOutput, error)
}

func getBucketLocation(c context.Context, api S3GetbucketLocation, bucketName string) (*s3.GetBucketLocationOutput, error) {
	in := &s3.GetBucketLocationInput{
		Bucket: aws.String(bucketName),
	}
	return api.GetBucketLocation(c, in)
}

type S3GetPublicAccessBlockAPI interface {
	GetPublicAccessBlock(ctx context.Context,
		params *s3.GetPublicAccessBlockInput,
		optFns ...func(*s3.Options)) (*s3.GetPublicAccessBlockOutput, error)
}

func getPublicAccessBlock(c context.Context, api S3GetPublicAccessBlockAPI, bucketName string) (*s3.GetPublicAccessBlockOutput, error) {
	in := &s3.GetPublicAccessBlockInput{
		Bucket: aws.String(bucketName),
	}
	return api.GetPublicAccessBlock(c, in)
}

type S3GetBucketEncryption interface {
	GetBucketEncryption(ctx context.Context,
		params *s3.GetBucketEncryptionInput,
		optFns ...func(*s3.Options)) (*s3.GetBucketEncryptionOutput, error)
}

func getBucketEncryption(c context.Context, api S3GetBucketEncryption, bucketName string) (*s3.GetBucketEncryptionOutput, error) {
	in := &s3.GetBucketEncryptionInput{
		Bucket: aws.String(bucketName),
	}
	return api.GetBucketEncryption(c, in)
}

type S3GetBucketVersioningAPI interface {
	GetBucketVersioning(ctx context.Context,
		params *s3.GetBucketVersioningInput,
		optFns ...func(*s3.Options)) (*s3.GetBucketVersioningOutput, error)
}

func getBucketVersioning(c context.Context, api S3GetBucketVersioningAPI, bucketName string) (*s3.GetBucketVersioningOutput, error) {
	in := &s3.GetBucketVersioningInput{
		Bucket: aws.String(bucketName),
	}
	return api.GetBucketVersioning(c, in)
}
