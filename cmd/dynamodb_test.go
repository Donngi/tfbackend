package cmd

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type mockDynamoDBCreateTableAPI func(ctx context.Context,
	params *dynamodb.CreateTableInput,
	optFns ...func(*dynamodb.Options)) (*dynamodb.CreateTableOutput, error)

func (m mockDynamoDBCreateTableAPI) CreateTable(ctx context.Context,
	params *dynamodb.CreateTableInput,
	optFns ...func(*dynamodb.Options)) (*dynamodb.CreateTableOutput, error) {

	return m(ctx, params, optFns...)
}

func createMockDynamoDBCreateTableAPIReturnInputTableName(t *testing.T) DynamoDBCreateTableAPI {
	return mockDynamoDBCreateTableAPI(func(ctx context.Context,
		params *dynamodb.CreateTableInput,
		optFns ...func(*dynamodb.Options)) (*dynamodb.CreateTableOutput, error) {

		out := dynamodb.CreateTableOutput{
			TableDescription: &types.TableDescription{
				TableName: params.TableName,
			},
		}
		return &out, nil
	})
}

func Test_createDynamoDBTable_Success(t *testing.T) {
	type args struct {
		tableName   string
		billingMode string
	}
	tests := []struct {
		name    string
		args    args
		api     func(t *testing.T) DynamoDBCreateTableAPI
		want    string
		wantErr bool
	}{
		// Given
		{
			name: "S01: billingMode=PAY_PER_REQUEST",
			args: args{
				tableName:   "TestTable",
				billingMode: "PAY_PER_REQUEST",
			},
			api:     createMockDynamoDBCreateTableAPIReturnInputTableName,
			want:    "TestTable",
			wantErr: false,
		},
		{
			name: "S02: billingMode=PROVISIONED",
			args: args{
				tableName:   "TestTable",
				billingMode: "PROVISIONED",
			},
			api:     createMockDynamoDBCreateTableAPIReturnInputTableName,
			want:    "TestTable",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// When
			got, err := createDynamoDBTable(context.Background(), tt.api(t), tt.args.tableName, tt.args.billingMode)

			// Then
			if (err != nil) != tt.wantErr {
				t.Errorf("createDynamoDBTable() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil || *got.TableDescription.TableName != tt.want {
				t.Errorf("createDynamoDBTable() got %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_createDynamoDBTable_Failure(t *testing.T) {
	type args struct {
		tableName   string
		billingMode string
	}
	tests := []struct {
		name    string
		args    args
		api     func(t *testing.T) DynamoDBCreateTableAPI
		want    *dynamodb.CreateTableOutput
		wantErr bool
	}{
		// Given
		{
			name: "F01: billingMode=invalid",
			args: args{
				tableName:   "TestTable",
				billingMode: "invalid💀",
			},
			api:     createMockDynamoDBCreateTableAPIReturnInputTableName,
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// When
			got, err := createDynamoDBTable(context.Background(), tt.api(t), tt.args.tableName, tt.args.billingMode)

			// Then
			if (err != nil) != tt.wantErr {
				t.Errorf("createDynamoDBTable() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("createDynamoDBTable() got %v, want %v", got, tt.want)
			}
		})
	}
}

type mockDynamoDBDescribeTableAPI func(ctx context.Context,
	params *dynamodb.DescribeTableInput,
	optFns ...func(*dynamodb.Options)) (*dynamodb.DescribeTableOutput, error)

func (m mockDynamoDBDescribeTableAPI) DescribeTable(ctx context.Context,
	params *dynamodb.DescribeTableInput,
	optFns ...func(*dynamodb.Options)) (*dynamodb.DescribeTableOutput, error) {

	return m(ctx, params, optFns...)
}

func createMockDynamoDBDescribeTableAPIReturnInputTableName(t *testing.T) DynamoDBDescribeTableAPI {
	return mockDynamoDBDescribeTableAPI(func(ctx context.Context,
		params *dynamodb.DescribeTableInput,
		optFns ...func(*dynamodb.Options)) (*dynamodb.DescribeTableOutput, error) {

		out := dynamodb.DescribeTableOutput{
			Table: &types.TableDescription{
				TableName: params.TableName,
			},
		}
		return &out, nil
	})
}

func createMockDynamoDBDescribeTableAPIReturnError(t *testing.T) DynamoDBDescribeTableAPI {
	return mockDynamoDBDescribeTableAPI(func(ctx context.Context,
		params *dynamodb.DescribeTableInput,
		optFns ...func(*dynamodb.Options)) (*dynamodb.DescribeTableOutput, error) {

		return nil, errors.New("some error")
	})
}

func Test_describeDynamoDBTable_Success(t *testing.T) {
	type args struct {
		tableName string
		api       func(t *testing.T) DynamoDBDescribeTableAPI
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// Given
		{
			name: "S01: Happy path",
			args: args{
				tableName: "HappyTable",
				api:       createMockDynamoDBDescribeTableAPIReturnInputTableName,
			},
			want:    "HappyTable",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// When
			got, err := describeDynamoDBTable(context.Background(), tt.args.api(t), tt.args.tableName)

			// Then
			if (err != nil) != tt.wantErr {
				t.Errorf("describeDynamoDBTable() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if *got.Table.TableName != tt.want {
				t.Errorf("describeDynamoDBTable() got.Table.TableName %v, want %v", got.Table.TableName, tt.want)
			}
		})
	}
}

func Test_describeDynamoDBTable_Failure(t *testing.T) {
	type args struct {
		tableName string
		api       func(t *testing.T) DynamoDBDescribeTableAPI
	}
	tests := []struct {
		name    string
		args    args
		want    *dynamodb.DescribeTableOutput
		wantErr bool
	}{
		// Given
		{
			name: "F01: Some error",
			args: args{
				tableName: "ErrorTable",
				api:       createMockDynamoDBDescribeTableAPIReturnError,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// When
			got, err := describeDynamoDBTable(context.Background(), tt.args.api(t), tt.args.tableName)

			// Then
			if (err != nil) != tt.wantErr {
				t.Errorf("describeDynamoDBTable() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("describeDynamoDBTable() got %v, want %v", got, tt.want)
			}
		})
	}
}
