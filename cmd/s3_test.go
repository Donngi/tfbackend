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
		name string
		args args
		api  func(t *testing.T) S3CreateBucketAPI
		want string
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
			want: "ap-northeast-1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// When
			got, err := createS3Bucket(context.Background(), tt.api(t), tt.args.bucketName, tt.args.region)

			// Then
			if err != nil {
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
		name string
		args args
		api  func(t *testing.T) S3CreateBucketAPI
		want string
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

					return nil, errors.New("Duplicated error")
				})
			},
			want: "failure💀",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// When
			_, err := createS3Bucket(context.Background(), tt.api(t), tt.args.bucketName, tt.args.region)

			// Then
			if err == nil {
				t.Errorf("createBucket() error = %v, want = some error", err)
			}
		})
	}
}
