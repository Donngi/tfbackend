package cmd

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// -----------------------------------
// For initS3 test
// -----------------------------------

func mockCreateBucketOK(ctx context.Context, params *s3.CreateBucketInput, optFns ...func(*s3.Options)) (*s3.CreateBucketOutput, error) {
	return &s3.CreateBucketOutput{}, nil
}
func mockCreateBucketNG(ctx context.Context, params *s3.CreateBucketInput, optFns ...func(*s3.Options)) (*s3.CreateBucketOutput, error) {
	return nil, errors.New("some error")
}

func mockPutPublicAccessBlockOK(ctx context.Context, params *s3.PutPublicAccessBlockInput, optFns ...func(*s3.Options)) (*s3.PutPublicAccessBlockOutput, error) {
	return &s3.PutPublicAccessBlockOutput{}, nil
}
func mockPutPublicAccessBlockNG(ctx context.Context, params *s3.PutPublicAccessBlockInput, optFns ...func(*s3.Options)) (*s3.PutPublicAccessBlockOutput, error) {
	return nil, errors.New("some error")
}

func mockPutBucketEncryptionOK(ctx context.Context, params *s3.PutBucketEncryptionInput, optFns ...func(*s3.Options)) (*s3.PutBucketEncryptionOutput, error) {
	return &s3.PutBucketEncryptionOutput{}, nil
}
func mockPutBucketEncryptionNG(ctx context.Context, params *s3.PutBucketEncryptionInput, optFns ...func(*s3.Options)) (*s3.PutBucketEncryptionOutput, error) {
	return nil, errors.New("some error")
}

func mockPutBucketVersioningOK(ctx context.Context, params *s3.PutBucketVersioningInput, optFns ...func(*s3.Options)) (*s3.PutBucketVersioningOutput, error) {
	return &s3.PutBucketVersioningOutput{}, nil
}
func mockPutBucketVersioningNG(ctx context.Context, params *s3.PutBucketVersioningInput, optFns ...func(*s3.Options)) (*s3.PutBucketVersioningOutput, error) {
	return nil, errors.New("some error")
}

func mockGetBucketLocationOK(ctx context.Context, params *s3.GetBucketLocationInput, optFns ...func(*s3.Options)) (*s3.GetBucketLocationOutput, error) {
	return &s3.GetBucketLocationOutput{
		LocationConstraint: "ap-northeast-1",
	}, nil
}
func mockGetBucketLocationNG(ctx context.Context, params *s3.GetBucketLocationInput, optFns ...func(*s3.Options)) (*s3.GetBucketLocationOutput, error) {
	return nil, errors.New("some error")
}

func mockGetPublicAccessBlockOK(ctx context.Context, params *s3.GetPublicAccessBlockInput, optFns ...func(*s3.Options)) (*s3.GetPublicAccessBlockOutput, error) {
	return &s3.GetPublicAccessBlockOutput{
		PublicAccessBlockConfiguration: &s3types.PublicAccessBlockConfiguration{
			BlockPublicAcls:       true,
			BlockPublicPolicy:     true,
			IgnorePublicAcls:      true,
			RestrictPublicBuckets: true,
		},
	}, nil
}
func mockGetPublicAccessBlockNG(ctx context.Context, params *s3.GetPublicAccessBlockInput, optFns ...func(*s3.Options)) (*s3.GetPublicAccessBlockOutput, error) {
	return nil, errors.New("some error")
}

func mockGetBucketEncryptionOK(ctx context.Context, params *s3.GetBucketEncryptionInput, optFns ...func(*s3.Options)) (*s3.GetBucketEncryptionOutput, error) {
	return &s3.GetBucketEncryptionOutput{
		ServerSideEncryptionConfiguration: &s3types.ServerSideEncryptionConfiguration{
			Rules: []s3types.ServerSideEncryptionRule{
				{
					ApplyServerSideEncryptionByDefault: &s3types.ServerSideEncryptionByDefault{
						SSEAlgorithm: s3types.ServerSideEncryptionAes256,
					},
				},
			},
		},
	}, nil
}
func mockGetBucketEncryptionNG(ctx context.Context, params *s3.GetBucketEncryptionInput, optFns ...func(*s3.Options)) (*s3.GetBucketEncryptionOutput, error) {
	return nil, errors.New("some error")
}

func mockGetBucketVersioningOK(ctx context.Context, params *s3.GetBucketVersioningInput, optFns ...func(*s3.Options)) (*s3.GetBucketVersioningOutput, error) {
	return &s3.GetBucketVersioningOutput{
		Status: s3types.BucketVersioningStatusEnabled,
	}, nil
}
func mockGetBucketVersioningNG(ctx context.Context, params *s3.GetBucketVersioningInput, optFns ...func(*s3.Options)) (*s3.GetBucketVersioningOutput, error) {
	return nil, errors.New("some error")
}

type mockS3ClientAllSuccess struct{}

func (m mockS3ClientAllSuccess) CreateBucket(ctx context.Context, params *s3.CreateBucketInput, optFns ...func(*s3.Options)) (*s3.CreateBucketOutput, error) {
	return mockCreateBucketOK(ctx, params, optFns...)
}
func (m mockS3ClientAllSuccess) PutPublicAccessBlock(ctx context.Context, params *s3.PutPublicAccessBlockInput, optFns ...func(*s3.Options)) (*s3.PutPublicAccessBlockOutput, error) {
	return mockPutPublicAccessBlockOK(ctx, params, optFns...)
}
func (m mockS3ClientAllSuccess) PutBucketEncryption(ctx context.Context, params *s3.PutBucketEncryptionInput, optFns ...func(*s3.Options)) (*s3.PutBucketEncryptionOutput, error) {
	return mockPutBucketEncryptionOK(ctx, params, optFns...)
}
func (m mockS3ClientAllSuccess) PutBucketVersioning(ctx context.Context, params *s3.PutBucketVersioningInput, optFns ...func(*s3.Options)) (*s3.PutBucketVersioningOutput, error) {
	return mockPutBucketVersioningOK(ctx, params, optFns...)
}
func (m mockS3ClientAllSuccess) GetBucketLocation(ctx context.Context, params *s3.GetBucketLocationInput, optFns ...func(*s3.Options)) (*s3.GetBucketLocationOutput, error) {
	return mockGetBucketLocationOK(ctx, params, optFns...)
}
func (m mockS3ClientAllSuccess) GetPublicAccessBlock(ctx context.Context, params *s3.GetPublicAccessBlockInput, optFns ...func(*s3.Options)) (*s3.GetPublicAccessBlockOutput, error) {
	return mockGetPublicAccessBlockOK(ctx, params, optFns...)
}
func (m mockS3ClientAllSuccess) GetBucketEncryption(ctx context.Context, params *s3.GetBucketEncryptionInput, optFns ...func(*s3.Options)) (*s3.GetBucketEncryptionOutput, error) {
	return mockGetBucketEncryptionOK(ctx, params, optFns...)
}
func (m mockS3ClientAllSuccess) GetBucketVersioning(ctx context.Context, params *s3.GetBucketVersioningInput, optFns ...func(*s3.Options)) (*s3.GetBucketVersioningOutput, error) {
	return mockGetBucketVersioningOK(ctx, params, optFns...)
}

type mockS3ClientCreateBucketFailure struct{}

func (m mockS3ClientCreateBucketFailure) CreateBucket(ctx context.Context, params *s3.CreateBucketInput, optFns ...func(*s3.Options)) (*s3.CreateBucketOutput, error) {
	return mockCreateBucketNG(ctx, params, optFns...)
}
func (m mockS3ClientCreateBucketFailure) PutPublicAccessBlock(ctx context.Context, params *s3.PutPublicAccessBlockInput, optFns ...func(*s3.Options)) (*s3.PutPublicAccessBlockOutput, error) {
	return mockPutPublicAccessBlockOK(ctx, params, optFns...)
}
func (m mockS3ClientCreateBucketFailure) PutBucketEncryption(ctx context.Context, params *s3.PutBucketEncryptionInput, optFns ...func(*s3.Options)) (*s3.PutBucketEncryptionOutput, error) {
	return mockPutBucketEncryptionOK(ctx, params, optFns...)
}
func (m mockS3ClientCreateBucketFailure) PutBucketVersioning(ctx context.Context, params *s3.PutBucketVersioningInput, optFns ...func(*s3.Options)) (*s3.PutBucketVersioningOutput, error) {
	return mockPutBucketVersioningOK(ctx, params, optFns...)
}
func (m mockS3ClientCreateBucketFailure) GetBucketLocation(ctx context.Context, params *s3.GetBucketLocationInput, optFns ...func(*s3.Options)) (*s3.GetBucketLocationOutput, error) {
	return mockGetBucketLocationOK(ctx, params, optFns...)
}
func (m mockS3ClientCreateBucketFailure) GetPublicAccessBlock(ctx context.Context, params *s3.GetPublicAccessBlockInput, optFns ...func(*s3.Options)) (*s3.GetPublicAccessBlockOutput, error) {
	return mockGetPublicAccessBlockOK(ctx, params, optFns...)
}
func (m mockS3ClientCreateBucketFailure) GetBucketEncryption(ctx context.Context, params *s3.GetBucketEncryptionInput, optFns ...func(*s3.Options)) (*s3.GetBucketEncryptionOutput, error) {
	return mockGetBucketEncryptionOK(ctx, params, optFns...)
}
func (m mockS3ClientCreateBucketFailure) GetBucketVersioning(ctx context.Context, params *s3.GetBucketVersioningInput, optFns ...func(*s3.Options)) (*s3.GetBucketVersioningOutput, error) {
	return mockGetBucketVersioningOK(ctx, params, optFns...)
}

type mockS3ClientPutPublicAccessBlockFailure struct{}

func (m mockS3ClientPutPublicAccessBlockFailure) CreateBucket(ctx context.Context, params *s3.CreateBucketInput, optFns ...func(*s3.Options)) (*s3.CreateBucketOutput, error) {
	return mockCreateBucketOK(ctx, params, optFns...)
}
func (m mockS3ClientPutPublicAccessBlockFailure) PutPublicAccessBlock(ctx context.Context, params *s3.PutPublicAccessBlockInput, optFns ...func(*s3.Options)) (*s3.PutPublicAccessBlockOutput, error) {
	return mockPutPublicAccessBlockNG(ctx, params, optFns...)
}
func (m mockS3ClientPutPublicAccessBlockFailure) PutBucketEncryption(ctx context.Context, params *s3.PutBucketEncryptionInput, optFns ...func(*s3.Options)) (*s3.PutBucketEncryptionOutput, error) {
	return mockPutBucketEncryptionOK(ctx, params, optFns...)
}
func (m mockS3ClientPutPublicAccessBlockFailure) PutBucketVersioning(ctx context.Context, params *s3.PutBucketVersioningInput, optFns ...func(*s3.Options)) (*s3.PutBucketVersioningOutput, error) {
	return mockPutBucketVersioningOK(ctx, params, optFns...)
}
func (m mockS3ClientPutPublicAccessBlockFailure) GetBucketLocation(ctx context.Context, params *s3.GetBucketLocationInput, optFns ...func(*s3.Options)) (*s3.GetBucketLocationOutput, error) {
	return mockGetBucketLocationOK(ctx, params, optFns...)
}
func (m mockS3ClientPutPublicAccessBlockFailure) GetPublicAccessBlock(ctx context.Context, params *s3.GetPublicAccessBlockInput, optFns ...func(*s3.Options)) (*s3.GetPublicAccessBlockOutput, error) {
	return mockGetPublicAccessBlockOK(ctx, params, optFns...)
}
func (m mockS3ClientPutPublicAccessBlockFailure) GetBucketEncryption(ctx context.Context, params *s3.GetBucketEncryptionInput, optFns ...func(*s3.Options)) (*s3.GetBucketEncryptionOutput, error) {
	return mockGetBucketEncryptionOK(ctx, params, optFns...)
}
func (m mockS3ClientPutPublicAccessBlockFailure) GetBucketVersioning(ctx context.Context, params *s3.GetBucketVersioningInput, optFns ...func(*s3.Options)) (*s3.GetBucketVersioningOutput, error) {
	return mockGetBucketVersioningOK(ctx, params, optFns...)
}

type mockS3ClientPutBucketEncryptionFailure struct{}

func (m mockS3ClientPutBucketEncryptionFailure) CreateBucket(ctx context.Context, params *s3.CreateBucketInput, optFns ...func(*s3.Options)) (*s3.CreateBucketOutput, error) {
	return mockCreateBucketOK(ctx, params, optFns...)
}
func (m mockS3ClientPutBucketEncryptionFailure) PutPublicAccessBlock(ctx context.Context, params *s3.PutPublicAccessBlockInput, optFns ...func(*s3.Options)) (*s3.PutPublicAccessBlockOutput, error) {
	return mockPutPublicAccessBlockOK(ctx, params, optFns...)
}
func (m mockS3ClientPutBucketEncryptionFailure) PutBucketEncryption(ctx context.Context, params *s3.PutBucketEncryptionInput, optFns ...func(*s3.Options)) (*s3.PutBucketEncryptionOutput, error) {
	return mockPutBucketEncryptionNG(ctx, params, optFns...)
}
func (m mockS3ClientPutBucketEncryptionFailure) PutBucketVersioning(ctx context.Context, params *s3.PutBucketVersioningInput, optFns ...func(*s3.Options)) (*s3.PutBucketVersioningOutput, error) {
	return mockPutBucketVersioningOK(ctx, params, optFns...)
}
func (m mockS3ClientPutBucketEncryptionFailure) GetBucketLocation(ctx context.Context, params *s3.GetBucketLocationInput, optFns ...func(*s3.Options)) (*s3.GetBucketLocationOutput, error) {
	return mockGetBucketLocationOK(ctx, params, optFns...)
}
func (m mockS3ClientPutBucketEncryptionFailure) GetPublicAccessBlock(ctx context.Context, params *s3.GetPublicAccessBlockInput, optFns ...func(*s3.Options)) (*s3.GetPublicAccessBlockOutput, error) {
	return mockGetPublicAccessBlockOK(ctx, params, optFns...)
}
func (m mockS3ClientPutBucketEncryptionFailure) GetBucketEncryption(ctx context.Context, params *s3.GetBucketEncryptionInput, optFns ...func(*s3.Options)) (*s3.GetBucketEncryptionOutput, error) {
	return mockGetBucketEncryptionOK(ctx, params, optFns...)
}
func (m mockS3ClientPutBucketEncryptionFailure) GetBucketVersioning(ctx context.Context, params *s3.GetBucketVersioningInput, optFns ...func(*s3.Options)) (*s3.GetBucketVersioningOutput, error) {
	return mockGetBucketVersioningOK(ctx, params, optFns...)
}

type mockS3ClientPutBucketVersioningFailure struct{}

func (m mockS3ClientPutBucketVersioningFailure) CreateBucket(ctx context.Context, params *s3.CreateBucketInput, optFns ...func(*s3.Options)) (*s3.CreateBucketOutput, error) {
	return mockCreateBucketOK(ctx, params, optFns...)
}
func (m mockS3ClientPutBucketVersioningFailure) PutPublicAccessBlock(ctx context.Context, params *s3.PutPublicAccessBlockInput, optFns ...func(*s3.Options)) (*s3.PutPublicAccessBlockOutput, error) {
	return mockPutPublicAccessBlockOK(ctx, params, optFns...)
}
func (m mockS3ClientPutBucketVersioningFailure) PutBucketEncryption(ctx context.Context, params *s3.PutBucketEncryptionInput, optFns ...func(*s3.Options)) (*s3.PutBucketEncryptionOutput, error) {
	return mockPutBucketEncryptionOK(ctx, params, optFns...)
}
func (m mockS3ClientPutBucketVersioningFailure) PutBucketVersioning(ctx context.Context, params *s3.PutBucketVersioningInput, optFns ...func(*s3.Options)) (*s3.PutBucketVersioningOutput, error) {
	return mockPutBucketVersioningNG(ctx, params, optFns...)
}
func (m mockS3ClientPutBucketVersioningFailure) GetBucketLocation(ctx context.Context, params *s3.GetBucketLocationInput, optFns ...func(*s3.Options)) (*s3.GetBucketLocationOutput, error) {
	return mockGetBucketLocationOK(ctx, params, optFns...)
}
func (m mockS3ClientPutBucketVersioningFailure) GetPublicAccessBlock(ctx context.Context, params *s3.GetPublicAccessBlockInput, optFns ...func(*s3.Options)) (*s3.GetPublicAccessBlockOutput, error) {
	return mockGetPublicAccessBlockOK(ctx, params, optFns...)
}
func (m mockS3ClientPutBucketVersioningFailure) GetBucketEncryption(ctx context.Context, params *s3.GetBucketEncryptionInput, optFns ...func(*s3.Options)) (*s3.GetBucketEncryptionOutput, error) {
	return mockGetBucketEncryptionOK(ctx, params, optFns...)
}
func (m mockS3ClientPutBucketVersioningFailure) GetBucketVersioning(ctx context.Context, params *s3.GetBucketVersioningInput, optFns ...func(*s3.Options)) (*s3.GetBucketVersioningOutput, error) {
	return mockGetBucketVersioningOK(ctx, params, optFns...)
}

type mockS3ClientGetBucketLocationFailure struct{}

func (m mockS3ClientGetBucketLocationFailure) CreateBucket(ctx context.Context, params *s3.CreateBucketInput, optFns ...func(*s3.Options)) (*s3.CreateBucketOutput, error) {
	return mockCreateBucketOK(ctx, params, optFns...)
}
func (m mockS3ClientGetBucketLocationFailure) PutPublicAccessBlock(ctx context.Context, params *s3.PutPublicAccessBlockInput, optFns ...func(*s3.Options)) (*s3.PutPublicAccessBlockOutput, error) {
	return mockPutPublicAccessBlockOK(ctx, params, optFns...)
}
func (m mockS3ClientGetBucketLocationFailure) PutBucketEncryption(ctx context.Context, params *s3.PutBucketEncryptionInput, optFns ...func(*s3.Options)) (*s3.PutBucketEncryptionOutput, error) {
	return mockPutBucketEncryptionOK(ctx, params, optFns...)
}
func (m mockS3ClientGetBucketLocationFailure) PutBucketVersioning(ctx context.Context, params *s3.PutBucketVersioningInput, optFns ...func(*s3.Options)) (*s3.PutBucketVersioningOutput, error) {
	return mockPutBucketVersioningOK(ctx, params, optFns...)
}
func (m mockS3ClientGetBucketLocationFailure) GetBucketLocation(ctx context.Context, params *s3.GetBucketLocationInput, optFns ...func(*s3.Options)) (*s3.GetBucketLocationOutput, error) {
	return mockGetBucketLocationNG(ctx, params, optFns...)
}
func (m mockS3ClientGetBucketLocationFailure) GetPublicAccessBlock(ctx context.Context, params *s3.GetPublicAccessBlockInput, optFns ...func(*s3.Options)) (*s3.GetPublicAccessBlockOutput, error) {
	return mockGetPublicAccessBlockOK(ctx, params, optFns...)
}
func (m mockS3ClientGetBucketLocationFailure) GetBucketEncryption(ctx context.Context, params *s3.GetBucketEncryptionInput, optFns ...func(*s3.Options)) (*s3.GetBucketEncryptionOutput, error) {
	return mockGetBucketEncryptionOK(ctx, params, optFns...)
}
func (m mockS3ClientGetBucketLocationFailure) GetBucketVersioning(ctx context.Context, params *s3.GetBucketVersioningInput, optFns ...func(*s3.Options)) (*s3.GetBucketVersioningOutput, error) {
	return mockGetBucketVersioningOK(ctx, params, optFns...)
}

type mockS3ClientGetPublicAccessBlockFailure struct{}

func (m mockS3ClientGetPublicAccessBlockFailure) CreateBucket(ctx context.Context, params *s3.CreateBucketInput, optFns ...func(*s3.Options)) (*s3.CreateBucketOutput, error) {
	return mockCreateBucketOK(ctx, params, optFns...)
}
func (m mockS3ClientGetPublicAccessBlockFailure) PutPublicAccessBlock(ctx context.Context, params *s3.PutPublicAccessBlockInput, optFns ...func(*s3.Options)) (*s3.PutPublicAccessBlockOutput, error) {
	return mockPutPublicAccessBlockOK(ctx, params, optFns...)
}
func (m mockS3ClientGetPublicAccessBlockFailure) PutBucketEncryption(ctx context.Context, params *s3.PutBucketEncryptionInput, optFns ...func(*s3.Options)) (*s3.PutBucketEncryptionOutput, error) {
	return mockPutBucketEncryptionOK(ctx, params, optFns...)
}
func (m mockS3ClientGetPublicAccessBlockFailure) PutBucketVersioning(ctx context.Context, params *s3.PutBucketVersioningInput, optFns ...func(*s3.Options)) (*s3.PutBucketVersioningOutput, error) {
	return mockPutBucketVersioningOK(ctx, params, optFns...)
}
func (m mockS3ClientGetPublicAccessBlockFailure) GetBucketLocation(ctx context.Context, params *s3.GetBucketLocationInput, optFns ...func(*s3.Options)) (*s3.GetBucketLocationOutput, error) {
	return mockGetBucketLocationOK(ctx, params, optFns...)
}
func (m mockS3ClientGetPublicAccessBlockFailure) GetPublicAccessBlock(ctx context.Context, params *s3.GetPublicAccessBlockInput, optFns ...func(*s3.Options)) (*s3.GetPublicAccessBlockOutput, error) {
	return mockGetPublicAccessBlockNG(ctx, params, optFns...)
}
func (m mockS3ClientGetPublicAccessBlockFailure) GetBucketEncryption(ctx context.Context, params *s3.GetBucketEncryptionInput, optFns ...func(*s3.Options)) (*s3.GetBucketEncryptionOutput, error) {
	return mockGetBucketEncryptionOK(ctx, params, optFns...)
}
func (m mockS3ClientGetPublicAccessBlockFailure) GetBucketVersioning(ctx context.Context, params *s3.GetBucketVersioningInput, optFns ...func(*s3.Options)) (*s3.GetBucketVersioningOutput, error) {
	return mockGetBucketVersioningOK(ctx, params, optFns...)
}

type mockS3ClientGetBucketEncryptionFailure struct{}

func (m mockS3ClientGetBucketEncryptionFailure) CreateBucket(ctx context.Context, params *s3.CreateBucketInput, optFns ...func(*s3.Options)) (*s3.CreateBucketOutput, error) {
	return mockCreateBucketOK(ctx, params, optFns...)
}
func (m mockS3ClientGetBucketEncryptionFailure) PutPublicAccessBlock(ctx context.Context, params *s3.PutPublicAccessBlockInput, optFns ...func(*s3.Options)) (*s3.PutPublicAccessBlockOutput, error) {
	return mockPutPublicAccessBlockOK(ctx, params, optFns...)
}
func (m mockS3ClientGetBucketEncryptionFailure) PutBucketEncryption(ctx context.Context, params *s3.PutBucketEncryptionInput, optFns ...func(*s3.Options)) (*s3.PutBucketEncryptionOutput, error) {
	return mockPutBucketEncryptionOK(ctx, params, optFns...)
}
func (m mockS3ClientGetBucketEncryptionFailure) PutBucketVersioning(ctx context.Context, params *s3.PutBucketVersioningInput, optFns ...func(*s3.Options)) (*s3.PutBucketVersioningOutput, error) {
	return mockPutBucketVersioningOK(ctx, params, optFns...)
}
func (m mockS3ClientGetBucketEncryptionFailure) GetBucketLocation(ctx context.Context, params *s3.GetBucketLocationInput, optFns ...func(*s3.Options)) (*s3.GetBucketLocationOutput, error) {
	return mockGetBucketLocationOK(ctx, params, optFns...)
}
func (m mockS3ClientGetBucketEncryptionFailure) GetPublicAccessBlock(ctx context.Context, params *s3.GetPublicAccessBlockInput, optFns ...func(*s3.Options)) (*s3.GetPublicAccessBlockOutput, error) {
	return mockGetPublicAccessBlockOK(ctx, params, optFns...)
}
func (m mockS3ClientGetBucketEncryptionFailure) GetBucketEncryption(ctx context.Context, params *s3.GetBucketEncryptionInput, optFns ...func(*s3.Options)) (*s3.GetBucketEncryptionOutput, error) {
	return mockGetBucketEncryptionNG(ctx, params, optFns...)
}
func (m mockS3ClientGetBucketEncryptionFailure) GetBucketVersioning(ctx context.Context, params *s3.GetBucketVersioningInput, optFns ...func(*s3.Options)) (*s3.GetBucketVersioningOutput, error) {
	return mockGetBucketVersioningOK(ctx, params, optFns...)
}

type mockS3ClientGetBucketVersioningFailure struct{}

func (m mockS3ClientGetBucketVersioningFailure) CreateBucket(ctx context.Context, params *s3.CreateBucketInput, optFns ...func(*s3.Options)) (*s3.CreateBucketOutput, error) {
	return mockCreateBucketOK(ctx, params, optFns...)
}
func (m mockS3ClientGetBucketVersioningFailure) PutPublicAccessBlock(ctx context.Context, params *s3.PutPublicAccessBlockInput, optFns ...func(*s3.Options)) (*s3.PutPublicAccessBlockOutput, error) {
	return mockPutPublicAccessBlockOK(ctx, params, optFns...)
}
func (m mockS3ClientGetBucketVersioningFailure) PutBucketEncryption(ctx context.Context, params *s3.PutBucketEncryptionInput, optFns ...func(*s3.Options)) (*s3.PutBucketEncryptionOutput, error) {
	return mockPutBucketEncryptionOK(ctx, params, optFns...)
}
func (m mockS3ClientGetBucketVersioningFailure) PutBucketVersioning(ctx context.Context, params *s3.PutBucketVersioningInput, optFns ...func(*s3.Options)) (*s3.PutBucketVersioningOutput, error) {
	return mockPutBucketVersioningOK(ctx, params, optFns...)
}
func (m mockS3ClientGetBucketVersioningFailure) GetBucketLocation(ctx context.Context, params *s3.GetBucketLocationInput, optFns ...func(*s3.Options)) (*s3.GetBucketLocationOutput, error) {
	return mockGetBucketLocationOK(ctx, params, optFns...)
}
func (m mockS3ClientGetBucketVersioningFailure) GetPublicAccessBlock(ctx context.Context, params *s3.GetPublicAccessBlockInput, optFns ...func(*s3.Options)) (*s3.GetPublicAccessBlockOutput, error) {
	return mockGetPublicAccessBlockOK(ctx, params, optFns...)
}
func (m mockS3ClientGetBucketVersioningFailure) GetBucketEncryption(ctx context.Context, params *s3.GetBucketEncryptionInput, optFns ...func(*s3.Options)) (*s3.GetBucketEncryptionOutput, error) {
	return mockGetBucketEncryptionOK(ctx, params, optFns...)
}
func (m mockS3ClientGetBucketVersioningFailure) GetBucketVersioning(ctx context.Context, params *s3.GetBucketVersioningInput, optFns ...func(*s3.Options)) (*s3.GetBucketVersioningOutput, error) {
	return mockGetBucketVersioningNG(ctx, params, optFns...)
}

// -----------------------------------
// For initDynamoDB test
// -----------------------------------

type mockDynamoDBClientAllSuccessProvisionedWrite5Read5 struct{}

func (m mockDynamoDBClientAllSuccessProvisionedWrite5Read5) CreateTable(ctx context.Context,
	params *dynamodb.CreateTableInput,
	optFns ...func(*dynamodb.Options)) (*dynamodb.CreateTableOutput, error) {
	return &dynamodb.CreateTableOutput{}, nil
}

func (m mockDynamoDBClientAllSuccessProvisionedWrite5Read5) DescribeTable(ctx context.Context,
	params *dynamodb.DescribeTableInput,
	optFns ...func(*dynamodb.Options)) (*dynamodb.DescribeTableOutput, error) {

	return &dynamodb.DescribeTableOutput{
		Table: &types.TableDescription{
			TableName: aws.String("happy-bucket"),
			BillingModeSummary: &types.BillingModeSummary{
				BillingMode: types.BillingModeProvisioned,
			},
			ProvisionedThroughput: &types.ProvisionedThroughputDescription{
				ReadCapacityUnits:  aws.Int64(5),
				WriteCapacityUnits: aws.Int64(5),
			},
		},
	}, nil
}

type mockDynamoDBClientAllSuccessProvisionedWrite5Read5WithoutBillingModeSummary struct{}

func (m mockDynamoDBClientAllSuccessProvisionedWrite5Read5WithoutBillingModeSummary) CreateTable(ctx context.Context,
	params *dynamodb.CreateTableInput,
	optFns ...func(*dynamodb.Options)) (*dynamodb.CreateTableOutput, error) {
	return &dynamodb.CreateTableOutput{}, nil
}

func (m mockDynamoDBClientAllSuccessProvisionedWrite5Read5WithoutBillingModeSummary) DescribeTable(ctx context.Context,
	params *dynamodb.DescribeTableInput,
	optFns ...func(*dynamodb.Options)) (*dynamodb.DescribeTableOutput, error) {

	return &dynamodb.DescribeTableOutput{
		Table: &types.TableDescription{
			TableName: aws.String("happy-bucket"),
			ProvisionedThroughput: &types.ProvisionedThroughputDescription{
				ReadCapacityUnits:  aws.Int64(5),
				WriteCapacityUnits: aws.Int64(5),
			},
		},
	}, nil
}

type mockDynamoDBClientAllSuccessPayPerRequest struct{}

func (m mockDynamoDBClientAllSuccessPayPerRequest) CreateTable(ctx context.Context,
	params *dynamodb.CreateTableInput,
	optFns ...func(*dynamodb.Options)) (*dynamodb.CreateTableOutput, error) {
	return &dynamodb.CreateTableOutput{}, nil
}

func (m mockDynamoDBClientAllSuccessPayPerRequest) DescribeTable(ctx context.Context,
	params *dynamodb.DescribeTableInput,
	optFns ...func(*dynamodb.Options)) (*dynamodb.DescribeTableOutput, error) {

	return &dynamodb.DescribeTableOutput{
		Table: &types.TableDescription{
			TableName: aws.String("happy-bucket"),
			BillingModeSummary: &types.BillingModeSummary{
				BillingMode: types.BillingModePayPerRequest,
			},
		},
	}, nil
}

type mockDynamoDBClientFailureCreateTableNG struct{}

func (m mockDynamoDBClientFailureCreateTableNG) CreateTable(ctx context.Context,
	params *dynamodb.CreateTableInput,
	optFns ...func(*dynamodb.Options)) (*dynamodb.CreateTableOutput, error) {
	return nil, errors.New("some error")
}

func (m mockDynamoDBClientFailureCreateTableNG) DescribeTable(ctx context.Context,
	params *dynamodb.DescribeTableInput,
	optFns ...func(*dynamodb.Options)) (*dynamodb.DescribeTableOutput, error) {

	return &dynamodb.DescribeTableOutput{}, nil
}

type mockDynamoDBClientFailureDescribeTableNG struct{}

func (m mockDynamoDBClientFailureDescribeTableNG) CreateTable(ctx context.Context,
	params *dynamodb.CreateTableInput,
	optFns ...func(*dynamodb.Options)) (*dynamodb.CreateTableOutput, error) {
	return &dynamodb.CreateTableOutput{}, nil
}

func (m mockDynamoDBClientFailureDescribeTableNG) DescribeTable(ctx context.Context,
	params *dynamodb.DescribeTableInput,
	optFns ...func(*dynamodb.Options)) (*dynamodb.DescribeTableOutput, error) {

	return nil, errors.New("some error")
}
