package cmd

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"unicode"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/olekukonko/tablewriter"
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

	GetBucketLocation(ctx context.Context,
		params *s3.GetBucketLocationInput,
		optFns ...func(*s3.Options)) (*s3.GetBucketLocationOutput, error)

	GetPublicAccessBlock(ctx context.Context,
		params *s3.GetPublicAccessBlockInput,
		optFns ...func(*s3.Options)) (*s3.GetPublicAccessBlockOutput, error)

	GetBucketEncryption(ctx context.Context,
		params *s3.GetBucketEncryptionInput,
		optFns ...func(*s3.Options)) (*s3.GetBucketEncryptionOutput, error)

	GetBucketVersioning(ctx context.Context,
		params *s3.GetBucketVersioningInput,
		optFns ...func(*s3.Options)) (*s3.GetBucketVersioningOutput, error)
}

type DynamoDBClientable interface {
	CreateTable(ctx context.Context,
		params *dynamodb.CreateTableInput,
		optFns ...func(*dynamodb.Options)) (*dynamodb.CreateTableOutput, error)
	DescribeTable(ctx context.Context,
		params *dynamodb.DescribeTableInput,
		optFns ...func(*dynamodb.Options)) (*dynamodb.DescribeTableOutput, error)
}

type initS3Result struct {
	BucketName        string
	Region            string
	BlockPublicAccess string
	Encryption        string
	Versioning        string
}

type initDynamoDBResult struct {
	TableName     string
	BillingMode   string
	WriteCapacity string
	ReadCapacity  string
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
	s3Res, err := initS3(s3, bucketName, region)
	if err != nil {
		return fmt.Errorf("failed to initialize s3 bucket: %w", err)
	}
	printCyan(fmt.Sprintf("Successfully create terraform backend - s3 bucket: %v\n", bucketName))
	fmt.Printf("Detail ... \n\n")
	sTable := tablewriter.NewWriter(os.Stdout)
	h, b := s3Res.createTableInput()
	sTable.SetHeader(h)
	for _, v := range b {
		sTable.Append(v)
	}
	sTable.SetAlignment(tablewriter.ALIGN_LEFT)
	sTable.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	sTable.SetCenterSeparator("|")
	sTable.Render()

	// Initialize DynamoDB table.
	if tableName != "" {
		dynamodb := dynamodb.NewFromConfig(cfg)
		dynamoRes, err := initDynamoDB(dynamodb, tableName, billingMode)
		if err != nil {
			return fmt.Errorf("failed to initialize dynamodb table: %w", err)
		}

		printCyan(fmt.Sprintf("Successfully create terraform lock table - dynamodb table: %v\n", tableName))
		fmt.Printf("Detail ... \n\n")
		dTable := tablewriter.NewWriter(os.Stdout)
		h, b := dynamoRes.createTableInput()
		dTable.SetHeader(h)
		for _, v := range b {
			dTable.Append(v)
		}
		dTable.SetAlignment(tablewriter.ALIGN_LEFT)
		dTable.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
		dTable.SetCenterSeparator("|")
		dTable.Render()
	}

	return nil
}

// initS3 setup terraform backend with messages.
func initS3(c S3Clientable, bucketName string, region string) (*initS3Result, error) {
	fmt.Printf("\n")
	fmt.Printf("---------------------------------------------------\n")
	fmt.Printf("🚀 Start to create terraform backend: s3 bucket ... \n")
	fmt.Printf("---------------------------------------------------\n")
	fmt.Printf("\n")

	// Create bucket
	fmt.Printf("Step1: Creating bucket ... ")
	if _, err := createS3Bucket(context.TODO(), c, bucketName, region); err != nil {
		printRed("FAILURE\n\n")
		return nil, fmt.Errorf("failed to create s3 bucket: %w", err)
	}
	fmt.Printf("SUCCESS\n")

	// Activate block all public access
	fmt.Printf("Step2: Activate block public access ... ")
	if _, err := enableAllPublicAccessBlock(context.TODO(), c, bucketName); err != nil {
		printRed("FAILURE\n\n")
		return nil, fmt.Errorf("failed to activate block public access of s3 bucket: %w", err)
	}
	fmt.Printf("SUCCESS\n")

	// Activate default encryption
	fmt.Printf("Step3: Activate default encryption (AES256) ... ")
	if _, err := enableBucketEncryptionAES256(context.TODO(), c, bucketName); err != nil {
		printRed("FAILURE\n\n")
		return nil, fmt.Errorf("failed to activate default encryption of s3 bucket: %w", err)
	}
	fmt.Printf("SUCCESS\n")

	// Activate versioning
	fmt.Printf("Step4: Activate bucket versioning ... ")
	if _, err := enableBucketVersioning(context.TODO(), c, bucketName); err != nil {
		printRed("FAILURE\n\n")
		return nil, fmt.Errorf("failed to activate versioning: %w", err)
	}
	fmt.Printf("SUCCESS\n")

	// Describe bucket
	res := initS3Result{
		BucketName: bucketName,
	}

	fmt.Printf("Step5: Confirmation - Get bucket location ... ")
	locationRes, err := getBucketLocation(context.TODO(), c, bucketName)
	if err != nil {
		printRed("FAILURE\n\n")
		return nil, fmt.Errorf("successfully created s3 bucket, but failed to describe s3 bucket: %w", err)
	}
	res.Region = string(locationRes.LocationConstraint)
	fmt.Printf("SUCCESS\n")

	fmt.Printf("Step6: Confirmation - Get block public access status ... ")
	blockRes, err := getPublicAccessBlock(context.TODO(), c, bucketName)
	if err != nil {
		printRed("FAILURE\n\n")
		return nil, fmt.Errorf("successfully created s3 bucket, but failed to describe s3 bucket: %w", err)
	}
	if blockRes.PublicAccessBlockConfiguration.BlockPublicAcls &&
		blockRes.PublicAccessBlockConfiguration.BlockPublicPolicy &&
		blockRes.PublicAccessBlockConfiguration.IgnorePublicAcls &&
		blockRes.PublicAccessBlockConfiguration.RestrictPublicBuckets {
		res.BlockPublicAccess = "Enabled"
	} else if !blockRes.PublicAccessBlockConfiguration.BlockPublicAcls &&
		!blockRes.PublicAccessBlockConfiguration.BlockPublicPolicy &&
		!blockRes.PublicAccessBlockConfiguration.IgnorePublicAcls &&
		!blockRes.PublicAccessBlockConfiguration.RestrictPublicBuckets {
		res.BlockPublicAccess = "Not fully enabled"
	}
	fmt.Printf("SUCCESS\n")

	fmt.Printf("Step7: Confirmation - Get bucket encryption status ... ")
	encryptionRes, err := getBucketEncryption(context.TODO(), c, bucketName)
	if err != nil {
		printRed("FAILURE\n\n")
		return nil, fmt.Errorf("successfully created s3 bucket, but failed to describe s3 bucket: %w", err)
	}
	if encryptionRes.ServerSideEncryptionConfiguration.Rules[0].ApplyServerSideEncryptionByDefault != nil {
		res.Encryption = string(encryptionRes.ServerSideEncryptionConfiguration.Rules[0].ApplyServerSideEncryptionByDefault.SSEAlgorithm)
	}
	fmt.Printf("SUCCESS\n")

	fmt.Printf("Step8: Confirmation - Get bucket versioning status ... ")
	versioningRes, err := getBucketVersioning(context.TODO(), c, bucketName)
	if err != nil {
		printRed("FAILURE\n\n")
		return nil, fmt.Errorf("successfully created s3 bucket, but failed to describe s3 bucket: %w", err)
	}
	res.Versioning = string(versioningRes.Status)
	fmt.Printf("SUCCESS\n")

	return &res, nil
}

// initDynamoDB setup terraform lock table with messages.
func initDynamoDB(c DynamoDBClientable, tableName string, billingMode string) (*initDynamoDBResult, error) {
	fmt.Printf("\n")
	fmt.Printf("---------------------------------------------------\n")
	fmt.Printf("🚀 Start to create terraform lock table: DynamoDB ... \n")
	fmt.Printf("---------------------------------------------------\n")
	fmt.Printf("\n")

	// Create table
	fmt.Printf("Step1: Creating table ... ")
	if _, err := createDynamoDBTable(context.TODO(), c, tableName, billingMode); err != nil {
		printRed("FAILURE\n\n")
		return nil, fmt.Errorf("failed to create dynamodb table: %w", err)
	}
	fmt.Printf("SUCCESS\n")

	// Describe table
	fmt.Printf("Step2: Confirmation - Describe table ... ")
	desc, err := describeDynamoDBTable(context.TODO(), c, tableName)
	if err != nil {
		printRed("FAILURE\n\n")
		return nil, fmt.Errorf("successfully created dynamodb table, but failed to describe dynamodb table: %w", err)
	}
	fmt.Printf("SUCCESS\n")

	res := initDynamoDBResult{}
	if desc.Table.TableName != nil {
		res.TableName = *desc.Table.TableName
	}

	if desc.Table.BillingModeSummary != nil {
		res.BillingMode = string(desc.Table.BillingModeSummary.BillingMode)
	} else if desc.Table.ProvisionedThroughput != nil {
		res.BillingMode = "PROVISIONED"
	}

	if desc.Table.ProvisionedThroughput != nil {
		res.WriteCapacity = strconv.FormatInt(*desc.Table.ProvisionedThroughput.WriteCapacityUnits, 10)
		res.ReadCapacity = strconv.FormatInt(*desc.Table.ProvisionedThroughput.ReadCapacityUnits, 10)
	}

	return &res, nil
}

func (i *initS3Result) createTableInput() (header []string, body [][]string) {
	h := []string{"PARAMETER", "VALUE"}
	b := [][]string{
		{"Bucket name", i.BucketName},
		{"Region", i.Region},
		{"Block Public Access", i.BlockPublicAccess},
		{"Encryption", i.Encryption},
		{"Versioning", i.Versioning},
	}
	return h, b
}

func (i *initDynamoDBResult) createTableInput() (header []string, body [][]string) {
	h := []string{"PARAMETER", "VALUE"}
	if i.BillingMode == "PAY_PER_REQUEST" {
		b := [][]string{
			{"Table name", i.TableName},
			{"Billing mode", i.BillingMode},
		}
		return h, b
	} else {
		b := [][]string{
			{"Table name", i.TableName},
			{"Billing mode", i.BillingMode},
			{"Write capacity", i.WriteCapacity},
			{"Read capacity", i.ReadCapacity},
		}
		return h, b
	}
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
