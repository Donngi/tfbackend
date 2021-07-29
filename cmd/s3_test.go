package cmd

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type mockS3CreateBucketAPI func(ctx context.Context,
	params *s3.CreateBucketInput,
	optFns ...func(*s3.Options)) (*s3.CreateBucketOutput, error)

func (m mockS3CreateBucketAPI) CreateBucket(ctx context.Context,
	params *s3.CreateBucketInput,
	optFns ...func(*s3.Options)) (*s3.CreateBucketOutput, error) {

	return m(ctx, params, optFns...)
}

func Test_createBucket_Success(t *testing.T) {
	type args struct {
		bucketName string
		region     string
	}

	// Given
	tests := []struct {
		name    string
		args    args
		api     func(t *testing.T) S3CreateBucketAPI
		want    string
		wantErr bool
	}{
		{
			name: "S01: Happy path",
			args: args{bucketName: "sample-bucket", region: "ap-northeast-1"},
			api: func(t *testing.T) S3CreateBucketAPI {
				return mockS3CreateBucketAPI(func(ctx context.Context,
					params *s3.CreateBucketInput,
					optFns ...func(*s3.Options)) (*s3.CreateBucketOutput, error) {

					output := &s3.CreateBucketOutput{
						Location: aws.String("ap-northeast-1"),
					}
					return output, nil
				})
			},
			want:    "ap-northeast-1",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// When
			got, err := createS3Bucket(context.Background(), tt.api(t), tt.args.bucketName, tt.args.region)

			// Then
			if (err != nil) != tt.wantErr {
				t.Errorf("createBucket() error = %v, want = nil", err)
			}
			if *got.Location != tt.want {
				t.Errorf("createBucket() got = %v, want = %v", got.Location, tt.want)
			}
		})
	}
}

func Test_createBucket_Failure(t *testing.T) {
	type args struct {
		bucketName string
		region     string
	}

	// Given
	tests := []struct {
		name    string
		args    args
		api     func(t *testing.T) S3CreateBucketAPI
		want    string
		wantErr bool
	}{
		{
			name: "F01: Duplicated bucket name",
			args: args{
				bucketName: "duplicated-bucket",
				region:     "ap-northeast-1",
			},
			api: func(t *testing.T) S3CreateBucketAPI {
				return mockS3CreateBucketAPI(func(ctx context.Context,
					params *s3.CreateBucketInput,
					optFns ...func(*s3.Options)) (*s3.CreateBucketOutput, error) {

					return nil, errors.New("duplicated error")
				})
			},
			want:    "failure💀",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// When
			_, err := createS3Bucket(context.Background(), tt.api(t), tt.args.bucketName, tt.args.region)

			// Then
			if (err != nil) != tt.wantErr {
				t.Errorf("createBucket() error = %v, want = some error", err)
			}
		})
	}
}

type mockS3PutPublicAccessBlockAPI func(ctx context.Context,
	params *s3.PutPublicAccessBlockInput,
	optFns ...func(*s3.Options)) (*s3.PutPublicAccessBlockOutput, error)

func (m mockS3PutPublicAccessBlockAPI) PutPublicAccessBlock(ctx context.Context,
	params *s3.PutPublicAccessBlockInput,
	optFns ...func(*s3.Options)) (*s3.PutPublicAccessBlockOutput, error) {

	return m(ctx, params, optFns...)
}

func Test_enableAllPublicAccessBlock_Success(t *testing.T) {
	type args struct {
		bucketName string
	}
	tests := []struct {
		name    string
		args    args
		api     func(t *testing.T) S3PutPublicAccessBlockAPI
		wantErr bool
	}{
		{
			name: "S01: Happy path",
			args: args{
				bucketName: "happy-bucket",
			},
			api: func(t *testing.T) S3PutPublicAccessBlockAPI {
				return mockS3PutPublicAccessBlockAPI(func(ctx context.Context,
					params *s3.PutPublicAccessBlockInput,
					optFns ...func(*s3.Options)) (*s3.PutPublicAccessBlockOutput, error) {

					return &s3.PutPublicAccessBlockOutput{}, nil
				})
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := enableAllPublicAccessBlock(context.Background(), tt.api(t), tt.args.bucketName)
			if (err != nil) != tt.wantErr {
				t.Errorf("enableAllPublicAccessBlock() error = %v, wantErr = %v", err, tt.wantErr)
				return
			}
			if got == nil {
				t.Errorf("enableAllPublicAccessBlock() got = nil, want some object")
			}
		})
	}
}

func Test_enableAllPublicAccessBlock_Failure(t *testing.T) {
	type args struct {
		bucketName string
	}
	tests := []struct {
		name    string
		args    args
		api     func(t *testing.T) S3PutPublicAccessBlockAPI
		wantErr bool
	}{
		{
			name: "F01: Bucket not exists",
			args: args{
				bucketName: "no-exist-bucket",
			},
			api: func(t *testing.T) S3PutPublicAccessBlockAPI {
				return mockS3PutPublicAccessBlockAPI(func(ctx context.Context,
					params *s3.PutPublicAccessBlockInput,
					optFns ...func(*s3.Options)) (*s3.PutPublicAccessBlockOutput, error) {

					return nil, errors.New("no exist error")
				})
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := enableAllPublicAccessBlock(context.Background(), tt.api(t), tt.args.bucketName)
			if (err != nil) != tt.wantErr {
				t.Errorf("enableAllPublicAccessBlock() error = %v, wantErr = %v", err, tt.wantErr)
				return
			}
		})
	}
}

type mockS3PutBucketEncryptionAPI func(ctx context.Context,
	params *s3.PutBucketEncryptionInput,
	optFns ...func(*s3.Options)) (*s3.PutBucketEncryptionOutput, error)

func (m mockS3PutBucketEncryptionAPI) PutBucketEncryption(ctx context.Context,
	params *s3.PutBucketEncryptionInput,
	optFns ...func(*s3.Options)) (*s3.PutBucketEncryptionOutput, error) {

	return m(ctx, params, optFns...)
}

func Test_enableBucketEncryptionAES256_Success(t *testing.T) {
	type args struct {
		bucketName string
	}
	tests := []struct {
		name    string
		args    args
		api     func(t *testing.T) mockS3PutBucketEncryptionAPI
		wantErr bool
	}{
		{
			name: "S01: Happy path",
			args: args{
				bucketName: "happy-bucket",
			},
			api: func(t *testing.T) mockS3PutBucketEncryptionAPI {
				return mockS3PutBucketEncryptionAPI(func(ctx context.Context,
					params *s3.PutBucketEncryptionInput,
					optFns ...func(*s3.Options)) (*s3.PutBucketEncryptionOutput, error) {

					return &s3.PutBucketEncryptionOutput{}, nil
				})
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := enableBucketEncryptionAES256(context.Background(), tt.api(t), tt.args.bucketName)
			if (err != nil) != tt.wantErr {
				t.Errorf("enableBucketEncryptionAES256() error = %v, wantErr = %v", err, tt.wantErr)
				return
			}
			if got == nil {
				t.Errorf("enableBucketEncryptionAES256() got = nil, want some object")
			}
		})
	}
}

func Test_enableBucketEncryptionAES256_Failure(t *testing.T) {
	type args struct {
		bucketName string
	}
	tests := []struct {
		name    string
		args    args
		api     func(t *testing.T) mockS3PutBucketEncryptionAPI
		wantErr bool
	}{
		{
			name: "F01: Bucket not exists",
			args: args{
				bucketName: "no-exist-bucket",
			},
			api: func(t *testing.T) mockS3PutBucketEncryptionAPI {
				return mockS3PutBucketEncryptionAPI(func(ctx context.Context,
					params *s3.PutBucketEncryptionInput,
					optFns ...func(*s3.Options)) (*s3.PutBucketEncryptionOutput, error) {

					return nil, errors.New("no exist error")
				})
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := enableBucketEncryptionAES256(context.Background(), tt.api(t), tt.args.bucketName)
			if (err != nil) != tt.wantErr {
				t.Errorf("enableBucketEncryptionAES256() error = %v, wantErr = %v", err, tt.wantErr)
				return
			}
		})
	}
}

type mockS3PutBucketVersioningAPI func(ctx context.Context,
	params *s3.PutBucketVersioningInput,
	optFns ...func(*s3.Options)) (*s3.PutBucketVersioningOutput, error)

func (m mockS3PutBucketVersioningAPI) PutBucketVersioning(ctx context.Context,
	params *s3.PutBucketVersioningInput,
	optFns ...func(*s3.Options)) (*s3.PutBucketVersioningOutput, error) {

	return m(ctx, params, optFns...)
}

func Test_enableBucketVersioning_Success(t *testing.T) {
	type args struct {
		bucketName string
	}
	tests := []struct {
		name    string
		args    args
		api     func(t *testing.T) S3PutBucketVersioningAPI
		wantErr bool
	}{
		{
			name: "S01: Happy path",
			args: args{
				bucketName: "happy-bucket",
			},
			api: func(t *testing.T) S3PutBucketVersioningAPI {
				return mockS3PutBucketVersioningAPI(func(ctx context.Context,
					params *s3.PutBucketVersioningInput,
					optFns ...func(*s3.Options)) (*s3.PutBucketVersioningOutput, error) {

					return &s3.PutBucketVersioningOutput{}, nil
				})
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := enableBucketVersioning(context.Background(), tt.api(t), tt.args.bucketName)
			if (err != nil) != tt.wantErr {
				t.Errorf("enableBucketVersioning() error = %v, wantErr = %v", err, tt.wantErr)
				return
			}
			if got == nil {
				t.Errorf("enableBucketVersioning() got = nil, want some object")
			}
		})
	}
}

func Test_enableBucketVersioning_Failure(t *testing.T) {
	type args struct {
		bucketName string
	}
	tests := []struct {
		name    string
		args    args
		api     func(t *testing.T) S3PutBucketVersioningAPI
		wantErr bool
	}{
		{
			name: "F01: Bucket not exists",
			args: args{
				bucketName: "no-exist-bucket",
			},
			api: func(t *testing.T) S3PutBucketVersioningAPI {
				return mockS3PutBucketVersioningAPI(func(ctx context.Context,
					params *s3.PutBucketVersioningInput,
					optFns ...func(*s3.Options)) (*s3.PutBucketVersioningOutput, error) {

					return nil, errors.New("no exist error")
				})
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := enableBucketVersioning(context.Background(), tt.api(t), tt.args.bucketName)
			if (err != nil) != tt.wantErr {
				t.Errorf("enableBucketVersioning() error = %v, wantErr = %v", err, tt.wantErr)
				return
			}
		})
	}
}
