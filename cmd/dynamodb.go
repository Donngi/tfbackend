package cmd

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type DynamoDBCreateTableAPI interface {
	CreateTable(ctx context.Context,
		params *dynamodb.CreateTableInput,
		optFns ...func(*dynamodb.Options)) (*dynamodb.CreateTableOutput, error)
}

func createDynamoDBTable(c context.Context, api DynamoDBCreateTableAPI, tableName string, billingMode string) (*dynamodb.CreateTableOutput, error) {
	if billingMode != string(types.BillingModePayPerRequest) && billingMode != string(types.BillingModeProvisioned) {
		return nil, fmt.Errorf("invalid billing mode")
	}

	in := &dynamodb.CreateTableInput{
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("LockID"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("LockID"),
				KeyType:       types.KeyTypeHash,
			},
		},
		TableName:   &tableName,
		BillingMode: types.BillingMode(billingMode),
	}

	if types.BillingMode(billingMode) == types.BillingModeProvisioned {
		in.ProvisionedThroughput = &types.ProvisionedThroughput{
			WriteCapacityUnits: aws.Int64(5),
			ReadCapacityUnits:  aws.Int64(5),
		}
	}

	return api.CreateTable(c, in)
}
