package cmd

import (
	"context"
	"fmt"
	"unicode"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/spf13/cobra"
)

var (
	bucketName  string
	region      string
	tableName   string
	billingMode string
)

type S3Clientable interface {
	CreateBucket(ctx context.Context,
		params *s3.CreateBucketInput,
		optFns ...func(*s3.Options)) (*s3.CreateBucketOutput, error)

	PutPublicAccessBlock(ctx context.Context,
		params *s3.PutPublicAccessBlockInput,
		optFns ...func(*s3.Options)) (*s3.PutPublicAccessBlockOutput, error)

	PutBucketEncryption(ctx context.Context,
		params *s3.PutBucketEncryptionInput,
		optFns ...func(*s3.Options)) (*s3.PutBucketEncryptionOutput, error)

	PutBucketVersioning(ctx context.Context,
		params *s3.PutBucketVersioningInput,
		optFns ...func(*s3.Options)) (*s3.PutBucketVersioningOutput, error)
}

type DynamoDBClientable interface {
	CreateTable(ctx context.Context,
		params *dynamodb.CreateTableInput,
		optFns ...func(*dynamodb.Options)) (*dynamodb.CreateTableOutput, error)
}

func NewCmdAws() *cobra.Command {
	cmd := &cobra.Command{

		Use:   "aws",
		Short: "Create S3 bucket and DynamoDB table for terraform backend and lock table.",
		Long: `Create S3 bucket and DynamoDB table for terraform backend and lock table.
	
By default, the bucket configuration is below.
- Enabled versioning
- Enabled block public access
- Enabled default encryption: SSE-S3(AES-256)

By default, the table configuration is below.
- Billing mode: PROVISIONED
`,
		SilenceUsage: true,
		RunE:         runCmdAws,
	}

	cobra.OnInitialize(func() {
		initBuildInDefault()
	})

	// flag
	cmd.Flags().StringVarP(&bucketName, "s3", "", "", "Name of S3 bucket to create.")
	cmd.MarkFlagRequired("s3")
	cmd.Flags().StringVarP(&tableName, "dynamodb", "", "", "Name of DynamoDB table to create.")
	cmd.Flags().StringVarP(&billingMode, "billing-mode", "", "", "DynamoDB billing mode. Only 'PAY_PER_REQUEST' or 'PROVISIONED' can be accepted. Default is PROVISIONED.")
	cmd.Flags().StringVarP(&region, "region", "", "", "Region where resources will be created. If not specified, tfbackend automatically set region from your cli configuration.")

	return cmd
}

func initBuildInDefault() {
	if billingMode == "" {
		billingMode = "PROVISIONED"
	}
}

func runCmdAws(cmd *cobra.Command, args []string) error {
	// Validation
	if !validateBucketName(bucketName) {
		return fmt.Errorf("bucket name contains capital letter: %v", bucketName)
	}

	// Load config
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return fmt.Errorf("configuration error: %w", err)
	}
	if region == "" {
		region = cfg.Region
	}

	// Initialize S3 bucket.
	s3 := s3.NewFromConfig(cfg)
	if err := initS3(s3, bucketName, region); err != nil {
		return fmt.Errorf("failed to initialize s3 bucket: %w", err)
	}

	// Initialize DynamoDB table.
	if tableName != "" {
		dynamodb := dynamodb.NewFromConfig(cfg)
		if err := initDynamoDB(dynamodb, tableName, billingMode); err != nil {
			return fmt.Errorf("failed to initialize dynamodb table: %w", err)
		}
	}

	return nil
}

// initS3 setup terraform backend with messages.
func initS3(c S3Clientable, bucketName string, region string) error {
	fmt.Printf("\n")
	fmt.Printf("Start to create terraform backend: s3 bucket ... \n")
	fmt.Printf("\n")

	// Create bucket
	fmt.Printf("Step1: Creating bucket ... ")
	if _, err := createS3Bucket(context.TODO(), c, bucketName, region); err != nil {
		printRed("FAILURE\n\n")
		return fmt.Errorf("failed to create s3 bucket: %w", err)
	}
	fmt.Printf("SUCCESS\n")

	// Activate block all public access
	fmt.Printf("Step2: Activate block public access ... ")
	if _, err := enableAllPublicAccessBlock(context.TODO(), c, bucketName); err != nil {
		printRed("FAILURE\n\n")
		return fmt.Errorf("failed to activate block public access of s3 bucket: %w", err)
	}
	fmt.Printf("SUCCESS\n")

	// Activate default encryption
	fmt.Printf("Step3: Activate default encryption (AES256) ... ")
	if _, err := enableBucketEncryptionAES256(context.TODO(), c, bucketName); err != nil {
		printRed("FAILURE\n\n")
		return fmt.Errorf("failed to activate default encryption of s3 bucket: %w", err)
	}
	fmt.Printf("SUCCESS\n")

	// Activate versioning
	fmt.Printf("Step4: Activate bucket versioning ... ")
	if _, err := enableBucketVersioning(context.TODO(), c, bucketName); err != nil {
		printRed("FAILURE\n\n")
		return fmt.Errorf("failed to activate versioning: %w", err)
	}
	fmt.Printf("SUCCESS\n")
	fmt.Printf("\n")

	printCyan(fmt.Sprintf("Successfully create terraform backend - s3 bucket: %v\n", bucketName))

	return nil
}

// initDynamoDB setup terraform lock table with messages.
func initDynamoDB(c DynamoDBClientable, tableName string, billingMode string) error {
	fmt.Printf("\n")
	fmt.Printf("Start to create terraform lock table: DynamoDB ... \n")
	fmt.Printf("\n")

	// Create table
	fmt.Printf("Step1: Creating table ... ")
	if _, err := createDynamoDBTable(context.TODO(), c, tableName, billingMode); err != nil {
		printRed("FAILURE\n\n")
		return fmt.Errorf("failed to create dynamodb table: %w", err)
	}
	fmt.Printf("SUCCESS\n")

	printCyan(fmt.Sprintf("Successfully create terraform lock table - dynamodb table: %v\n", tableName))

	return nil
}

// validateBucketName checks if bucket name contains capitals.
func validateBucketName(b string) bool {
	for _, r := range b {
		if !unicode.IsLower(r) && unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

func validateBillingMode(m string) bool {
	if m != string(types.BillingModePayPerRequest) && m != string(types.BillingModeProvisioned) {
		return false
	}
	return true
}
