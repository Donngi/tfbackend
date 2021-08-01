package cmd

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
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

type mockS3GetBucketLocation func(ctx context.Context, params *s3.GetBucketLocationInput, optFns ...func(*s3.Options)) (*s3.GetBucketLocationOutput, error)

func (m mockS3GetBucketLocation) GetBucketLocation(ctx context.Context, params *s3.GetBucketLocationInput, optFns ...func(*s3.Options)) (*s3.GetBucketLocationOutput, error) {
	return m(ctx, params, optFns...)
}

func TestGetBucketLocation(t *testing.T) {
	type args struct {
		api        func(t *testing.T) S3GetbucketLocation
		bucketName string
	}
	tests := []struct {
		name    string
		args    args
		want    *string
		wantErr bool
	}{
		{
			name: "S01: Happy path",
			args: args{
				bucketName: "happy-bucket",
				api: func(t *testing.T) S3GetbucketLocation {
					return mockS3GetBucketLocation(func(ctx context.Context, params *s3.GetBucketLocationInput, optFns ...func(*s3.Options)) (*s3.GetBucketLocationOutput, error) {
						out := &s3.GetBucketLocationOutput{
							LocationConstraint: types.BucketLocationConstraint("ap-northeast-1"),
						}
						return out, nil
					})
				},
			},
			want:    aws.String("ap-northeast-1"),
			wantErr: false,
		},
		{
			name: "F01: some error",
			args: args{
				bucketName: "error-bucket",
				api: func(t *testing.T) S3GetbucketLocation {
					return mockS3GetBucketLocation(func(ctx context.Context, params *s3.GetBucketLocationInput, optFns ...func(*s3.Options)) (*s3.GetBucketLocationOutput, error) {
						return nil, errors.New("some error")
					})
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetBucketLocation(context.Background(), tt.args.api(t), tt.args.bucketName)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBucketLocation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.want != nil {
				if got.LocationConstraint != types.BucketLocationConstraint(*tt.want) {
					t.Errorf("GetBucketLocation() got.LocationConstraint %v, want %v", got, tt.want)
				}
			}
		})
	}
}

type mockS3GetPublicAccessBlockAPI func(ctx context.Context,
	params *s3.GetPublicAccessBlockInput,
	optFns ...func(*s3.Options)) (*s3.GetPublicAccessBlockOutput, error)

func (m mockS3GetPublicAccessBlockAPI) GetPublicAccessBlock(ctx context.Context,
	params *s3.GetPublicAccessBlockInput,
	optFns ...func(*s3.Options)) (*s3.GetPublicAccessBlockOutput, error) {

	return m(ctx, params, optFns...)
}
func Test_getPublicAccessBlock(t *testing.T) {
	type args struct {
		bucketName string
		api        func(t *testing.T) S3GetPublicAccessBlockAPI
	}
	tests := []struct {
		name    string
		args    args
		want    *types.PublicAccessBlockConfiguration
		wantErr bool
	}{
		{
			name: "S01: Happy path",
			args: args{
				bucketName: "happy-bucket",
				api: func(t *testing.T) S3GetPublicAccessBlockAPI {
					return mockS3GetPublicAccessBlockAPI(func(ctx context.Context,
						params *s3.GetPublicAccessBlockInput,
						optFns ...func(*s3.Options)) (*s3.GetPublicAccessBlockOutput, error) {

						out := &s3.GetPublicAccessBlockOutput{
							PublicAccessBlockConfiguration: &types.PublicAccessBlockConfiguration{
								BlockPublicAcls:       true,
								BlockPublicPolicy:     true,
								IgnorePublicAcls:      true,
								RestrictPublicBuckets: true,
							},
						}
						return out, nil
					})
				},
			},
			want: &types.PublicAccessBlockConfiguration{
				BlockPublicAcls:       true,
				BlockPublicPolicy:     true,
				IgnorePublicAcls:      true,
				RestrictPublicBuckets: true,
			},
			wantErr: false,
		},
		{
			name: "F01: Some error",
			args: args{
				bucketName: "failure-bucket",
				api: func(t *testing.T) S3GetPublicAccessBlockAPI {
					return mockS3GetPublicAccessBlockAPI(func(ctx context.Context,
						params *s3.GetPublicAccessBlockInput,
						optFns ...func(*s3.Options)) (*s3.GetPublicAccessBlockOutput, error) {

						return nil, errors.New("some error")
					})
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getPublicAccessBlock(context.Background(), tt.args.api(t), tt.args.bucketName)
			if (err != nil) != tt.wantErr {
				t.Errorf("getPublicAccessBlock() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.want != nil {
				if !reflect.DeepEqual(got.PublicAccessBlockConfiguration, tt.want) {
					t.Errorf("getPublicAccessBlock() got.PublicAccessBlockConfiguration %v, want %v", got, tt.want)
				}
			}
		})
	}
}

type mockS3GetBucketEncryptionAPI func(ctx context.Context,
	params *s3.GetBucketEncryptionInput,
	optFns ...func(*s3.Options)) (*s3.GetBucketEncryptionOutput, error)

func (m mockS3GetBucketEncryptionAPI) GetBucketEncryption(ctx context.Context,
	params *s3.GetBucketEncryptionInput,
	optFns ...func(*s3.Options)) (*s3.GetBucketEncryptionOutput, error) {
	return m(ctx, params, optFns...)
}
func Test_getBucketEncryption(t *testing.T) {
	type args struct {
		api    func(t *testing.T) S3GetBucketEncryption
		optFns []func(*s3.Options)
	}
	tests := []struct {
		name    string
		args    args
		want    *types.ServerSideEncryptionConfiguration
		wantErr bool
	}{
		{
			name: "S01: Happy path",
			args: args{
				api: func(t *testing.T) S3GetBucketEncryption {
					return mockS3GetBucketEncryptionAPI(func(ctx context.Context, params *s3.GetBucketEncryptionInput, optFns ...func(*s3.Options)) (*s3.GetBucketEncryptionOutput, error) {
						out := &s3.GetBucketEncryptionOutput{
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
						return out, nil
					})
				},
			},
			want: &types.ServerSideEncryptionConfiguration{
				Rules: []types.ServerSideEncryptionRule{
					{
						ApplyServerSideEncryptionByDefault: &types.ServerSideEncryptionByDefault{
							SSEAlgorithm: types.ServerSideEncryptionAes256,
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "F01: Some error",
			args: args{
				api: func(t *testing.T) S3GetBucketEncryption {
					return mockS3GetBucketEncryptionAPI(func(ctx context.Context, params *s3.GetBucketEncryptionInput, optFns ...func(*s3.Options)) (*s3.GetBucketEncryptionOutput, error) {
						return nil, errors.New("some error")
					})
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getBucketEncryption(context.Background(), tt.args.api(t), tt.args.optFns...)
			if (err != nil) != tt.wantErr {
				t.Errorf("getBucketEncryption() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.want != nil {
				if !reflect.DeepEqual(got.ServerSideEncryptionConfiguration, tt.want) {
					t.Errorf("getBucketEncryption() got.ServerSideEncryptionConfiguration %v, want %v", got, tt.want)
				}
			}
		})
	}
}

type mockS3GetBucketVersioning func(ctx context.Context,
	params *s3.GetBucketVersioningInput,
	optFns ...func(*s3.Options)) (*s3.GetBucketVersioningOutput, error)

func (m mockS3GetBucketVersioning) GetBucketVersioning(ctx context.Context,
	params *s3.GetBucketVersioningInput,
	optFns ...func(*s3.Options)) (*s3.GetBucketVersioningOutput, error) {

	return m(ctx, params, optFns...)
}
func Test_getBucketVersioning(t *testing.T) {
	type args struct {
		api        func(t *testing.T) S3GetBucketVersioningAPI
		bucketName string
	}
	tests := []struct {
		name    string
		args    args
		want    types.BucketVersioningStatus
		wantErr bool
	}{
		{
			name: "S01: Happy path",
			args: args{
				bucketName: "happy-bucket",
				api: func(t *testing.T) S3GetBucketVersioningAPI {
					return mockS3GetBucketVersioning(func(ctx context.Context,
						params *s3.GetBucketVersioningInput,
						optFns ...func(*s3.Options)) (*s3.GetBucketVersioningOutput, error) {

						out := &s3.GetBucketVersioningOutput{
							Status: types.BucketVersioningStatusEnabled,
						}
						return out, nil
					})
				},
			},
			want:    types.BucketVersioningStatusEnabled,
			wantErr: false,
		},
		{
			name: "F01: Some error",
			args: args{
				bucketName: "failure-bucket",
				api: func(t *testing.T) S3GetBucketVersioningAPI {
					return mockS3GetBucketVersioning(func(ctx context.Context,
						params *s3.GetBucketVersioningInput,
						optFns ...func(*s3.Options)) (*s3.GetBucketVersioningOutput, error) {

						return nil, errors.New("some error")
					})
				},
			},
			want:    "failure💀",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getBucketVersioning(context.Background(), tt.args.api(t), tt.args.bucketName)
			if (err != nil) != tt.wantErr {
				t.Errorf("getBucketVersioning() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.want == types.BucketVersioningStatusEnabled || tt.want == types.BucketVersioningStatusSuspended {
				if !reflect.DeepEqual(got.Status, tt.want) {
					t.Errorf("getBucketVersioning() got.Status %v, want %v", got.Status, tt.want)
				}
			}
		})
	}
}
