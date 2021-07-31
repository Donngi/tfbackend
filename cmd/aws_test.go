package cmd

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func Test_validateBucketName(t *testing.T) {
	type args struct {
		bucketName string
	}

	// Given
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "S01: Happy path",
			args: args{bucketName: "success-bucket"},
			want: true,
		},
		{
			name: "F01: Contains capitals",
			args: args{bucketName: "FAILURE-BUCKET"},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// When
			res := validateBucketName(tt.args.bucketName)

			// Then
			if res != tt.want {
				t.Errorf("validateBucketName() res = %v, want = %v", res, tt.want)
			}
		})
	}
}

func Test_validateBillingMode(t *testing.T) {
	type args struct {
		mode string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "S01: PAY_PER_REQUEST",
			args: args{
				mode: "PAY_PER_REQUEST",
			},
			want: true,
		},
		{
			name: "S02: PROVISIONED",
			args: args{
				mode: "PROVISIONED",
			},
			want: true,
		},
		{
			name: "F01: invalid",
			args: args{
				mode: "invalid💀",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := validateBillingMode(tt.args.mode); got != tt.want {
				t.Errorf("validateBillingMode() = %v, want %v", got, tt.want)
			}
		})
	}
}

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

func Test_initDynamoDB(t *testing.T) {
	type args struct {
		c           DynamoDBClientable
		tableName   string
		billingMode string
	}
	tests := []struct {
		name    string
		args    args
		want    *initDynamoDBResult
		wantErr bool
	}{
		{
			name: "S01: PROVISIONED, Write=5, Read=5",
			args: args{
				c:           mockDynamoDBClientAllSuccessProvisionedWrite5Read5{},
				tableName:   "happy-bucket",
				billingMode: "PROVISIONED",
			},
			want: &initDynamoDBResult{
				TableName:     "happy-bucket",
				BillingMode:   "PROVISIONED",
				WriteCapacity: "5",
				ReadCapacity:  "5",
			},
			wantErr: false,
		},
		{
			name: "S02: PROVISIONED, Write=5, Read=5. AWS doesn't specify BillingModeSummary.",
			args: args{
				c:           mockDynamoDBClientAllSuccessProvisionedWrite5Read5WithoutBillingModeSummary{},
				tableName:   "happy-bucket",
				billingMode: "PROVISIONED",
			},
			want: &initDynamoDBResult{
				TableName:     "happy-bucket",
				BillingMode:   "PROVISIONED",
				WriteCapacity: "5",
				ReadCapacity:  "5",
			},
			wantErr: false,
		},
		{
			name: "S03: PAY_PER_REQUEST",
			args: args{
				c:           mockDynamoDBClientAllSuccessPayPerRequest{},
				tableName:   "happy-bucket",
				billingMode: "PAY_PER_REQUEST",
			},
			want: &initDynamoDBResult{
				TableName:   "happy-bucket",
				BillingMode: "PAY_PER_REQUEST",
			},
			wantErr: false,
		},
		{
			name: "F01: CreateTable fails",
			args: args{
				c:           mockDynamoDBClientFailureCreateTableNG{},
				tableName:   "failure-bucket",
				billingMode: "PAY_PER_REQUEST",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "F02: DescribeTable fails",
			args: args{
				c:           mockDynamoDBClientFailureDescribeTableNG{},
				tableName:   "failure-bucket",
				billingMode: "PAY_PER_REQUEST",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := initDynamoDB(tt.args.c, tt.args.tableName, tt.args.billingMode)
			if (err != nil) != tt.wantErr {
				t.Errorf("initDynamoDB() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("initDynamoDB() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_initDynamoDBResult_createTableInput(t *testing.T) {
	type fields struct {
		TableName     string
		BillingMode   string
		WriteCapacity string
		ReadCapacity  string
	}
	tests := []struct {
		name       string
		fields     fields
		wantHeader []string
		wantBody   [][]string
	}{
		{
			name: "S01: PROVISIONED",
			fields: fields{
				TableName:     "happy-table",
				BillingMode:   "PROVISIONED",
				WriteCapacity: "5",
				ReadCapacity:  "5",
			},
			wantHeader: []string{"PARAMETER", "VALUE"},
			wantBody: [][]string{
				{"Table name", "happy-table"},
				{"Billing mode", "PROVISIONED"},
				{"Write capacity", "5"},
				{"Read capacity", "5"},
			},
		},
		{
			name: "S02: PAY_PER_REQUEST",
			fields: fields{
				TableName:   "happy-table",
				BillingMode: "PAY_PER_REQUEST",
			},
			wantHeader: []string{"PARAMETER", "VALUE"},
			wantBody: [][]string{
				{"Table name", "happy-table"},
				{"Billing mode", "PAY_PER_REQUEST"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &initDynamoDBResult{
				TableName:     tt.fields.TableName,
				BillingMode:   tt.fields.BillingMode,
				WriteCapacity: tt.fields.WriteCapacity,
				ReadCapacity:  tt.fields.ReadCapacity,
			}
			gotHeader, gotBody := i.createTableInput()
			if !reflect.DeepEqual(gotHeader, tt.wantHeader) {
				t.Errorf("initDynamoDBResult.createTableInput() gotHeader = %v, want %v", gotHeader, tt.wantHeader)
			}
			if !reflect.DeepEqual(gotBody, tt.wantBody) {
				t.Errorf("initDynamoDBResult.createTableInput() gotBody = %v, want %v", gotBody, tt.wantBody)
			}
		})
	}
}
