package cmd

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type mockCreateBucketAPI struct{}

func (m mockCreateBucketAPI) CreateBucket(ctx context.Context,
	params *s3.CreateBucketInput,
	optFns ...func(*s3.Options)) (*s3.CreateBucketOutput, error) {

	output := &s3.CreateBucketOutput{
		Location: aws.String("ap-northeast-1"),
	}
	return output, nil
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
		want string
	}{
		{
			name: "S01: Happy path",
			args: args{bucketName: "sample-bucket", region: "ap-northeast-1"},
			want: "ap-northeast-1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// When
			api := mockCreateBucketAPI{}
			res, err := createS3Bucket(context.Background(), api, tt.args.bucketName, tt.args.region)

			// Then
			if err != nil {
				t.Errorf("createBucket() error = %v, want = nil", err)
			}
			if *res.Location != tt.want {
				t.Errorf("createBucket() res = %v, want = %v", res.Location, tt.want)
			}
		})
	}
}

// TODO
// func Test_createBucket_Failure(t *testing.T) {
// 	type args struct {
// 		bucketName string
// 	}

// 	// Given
// 	tests := []struct {
// 		name string
// 		args args
// 		want string
// 	}{
// 		{
// 			name: "F01: Duplicated bucket name",
// 			args: args{bucketName: "duplicated-bucket"},
// 			want: "failure💀",
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			// When
// 			api := mockCreateBucketAPI{}
// 			_, err := createS3Bucket(context.Background(), api, tt.args.bucketName)

// 			// Then
// 			if err == nil {
// 				t.Errorf("createBucket() error = %v, want = some error", err)
// 			}
// 		})
// 	}
// }

