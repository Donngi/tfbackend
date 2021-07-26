package cmd

import (
	"context"
	"fmt"
	"unicode"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/spf13/cobra"
)

var (
	bucketName string
	region     string
)

func NewCmdAws() *cobra.Command {
	cmd := &cobra.Command{

		Use:   "aws",
		Short: "Create s3 bucket for terraform backend.",
		Long: `Create s3 bucket for terraform backend.
	
By default, the bucket configuration is below.
- Enabled versioning
- Enabled block public access
- Enabled default encryption: SSE-S3(AES-256)`,
		SilenceUsage: true,
		RunE:         runCmdAws,
	}

	// flag
	cmd.Flags().StringVarP(&bucketName, "bucket", "", "", "S3 bucket to create")
	cmd.MarkFlagRequired("bucket")
	cmd.Flags().StringVarP(&bucketName, "region", "", "", "Region where S3 bucket will be created. If not specified, tfbackend automatically set region from your cli configuration.")

	return cmd
}

func runCmdAws(cmd *cobra.Command, args []string) error {
	// Validation
	if !validateBucketName(bucketName) {
		return fmt.Errorf("bucket name contains capital letter: %v", bucketName)
	}

	// Create bucket
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return fmt.Errorf("configuration error: %w", err)
	}
	if region == "" {
		region = cfg.Region
	}

	client := s3.NewFromConfig(cfg)

	fmt.Printf("\n")
	fmt.Printf("Start to create terraform backend s3 bucket ... \n")
	fmt.Printf("\n")

	fmt.Printf("Step1: Creating bucket ... ")
	if _, err := createS3Bucket(context.TODO(), client, bucketName, region); err != nil {
		printRed("FAILURE\n\n")
		return fmt.Errorf("failed to create s3 bucket: %w", err)
	}
	fmt.Printf("SUCCESS\n")

	// Activate block all public access
	fmt.Printf("Step2: Activate block public access ... ")
	if _, err := enableAllPublicAccessBlock(context.TODO(), client, bucketName); err != nil {
		printRed("FAILURE\n\n")
		return fmt.Errorf("failed to activate block public access of s3 bucket: %w", err)
	}
	fmt.Printf("SUCCESS\n")

	// Activate default encryption
	fmt.Printf("Step3: Activate default encryption (AES256) ... ")
	if _, err := enableBucketEncryptionAES256(context.TODO(), client, bucketName); err != nil {
		printRed("FAILURE\n\n")
		return fmt.Errorf("failed to activate default encryption of s3 bucket: %w", err)
	}
	fmt.Printf("SUCCESS\n")

	// Activate versioning
	fmt.Printf("Step4: Activate bucket versioning ... ")
	if _, err := enableBucketVersioning(context.TODO(), client, bucketName); err != nil {
		printRed("FAILURE\n\n")
		return fmt.Errorf("failed to activate versioning: %w", err)
	}
	fmt.Printf("SUCCESS\n")
	fmt.Printf("\n")

	printCyan(fmt.Sprintf("Successfully create terraform backend s3 bucket: %v\n", bucketName))

	return nil
}

func validateBucketName(b string) bool {
	// Check if bucket name contains capitals.
	for _, r := range b {
		if !unicode.IsLower(r) && unicode.IsLetter(r) {
			return false
		}
	}
	return true
}
