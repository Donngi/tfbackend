package cmd

import (
	"reflect"
	"testing"
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

func Test_initS3(t *testing.T) {
	type args struct {
		c          S3Clientable
		bucketName string
		region     string
	}
	tests := []struct {
		name    string
		args    args
		want    *initS3Result
		wantErr bool
	}{
		{
			name: "S01: Happy path",
			args: args{
				c:          mockS3ClientAllSuccess{},
				bucketName: "happy-bucket",
				region:     "ap-northeast-1",
			},
			want: &initS3Result{
				BucketName:        "happy-bucket",
				Region:            "ap-northeast-1",
				BlockPublicAccess: "Enabled",
				Encryption:        "AES256",
				Versioning:        "Enabled",
			},
			wantErr: false,
		},
		{
			name: "F01: CreateBucket fails",
			args: args{
				c:          mockS3ClientCreateBucketFailure{},
				bucketName: "error-bucket",
				region:     "ap-northeast-1",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "F02: PutPublicAccessBlock fails",
			args: args{
				c:          mockS3ClientPutPublicAccessBlockFailure{},
				bucketName: "error-bucket",
				region:     "ap-northeast-1",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "F03: PutBucketEncryption fails",
			args: args{
				c:          mockS3ClientPutBucketEncryptionFailure{},
				bucketName: "error-bucket",
				region:     "ap-northeast-1",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "F04: PutBucketVersioning fails",
			args: args{
				c:          mockS3ClientPutBucketVersioningFailure{},
				bucketName: "error-bucket",
				region:     "ap-northeast-1",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "F05: GetBucketLocation fails",
			args: args{
				c:          mockS3ClientGetBucketLocationFailure{},
				bucketName: "error-bucket",
				region:     "ap-northeast-1",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "F06: GetPublicAccessBlock fails",
			args: args{
				c:          mockS3ClientGetPublicAccessBlockFailure{},
				bucketName: "error-bucket",
				region:     "ap-northeast-1",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "F07: GetBucketEncryption fails",
			args: args{
				c:          mockS3ClientGetBucketEncryptionFailure{},
				bucketName: "error-bucket",
				region:     "ap-northeast-1",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "F08: GetBucketVersioning fails",
			args: args{
				c:          mockS3ClientGetBucketVersioningFailure{},
				bucketName: "error-bucket",
				region:     "ap-northeast-1",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := initS3(tt.args.c, tt.args.bucketName, tt.args.region)
			if (err != nil) != tt.wantErr {
				t.Errorf("initS3() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("initS3() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_initS3Result_createTableInput(t *testing.T) {
	type fields struct {
		BucketName        string
		Region            string
		BlockPublicAccess string
		Encryption        string
		Versioning        string
	}
	tests := []struct {
		name       string
		fields     fields
		wantHeader []string
		wantBody   [][]string
	}{
		{
			name: "S01: Happy path",
			fields: fields{
				BucketName:        "happy-bucket",
				Region:            "ap-northeast-1",
				BlockPublicAccess: "Enabled",
				Encryption:        "AES256",
				Versioning:        "Enabled",
			},
			wantHeader: []string{"PARAMETER", "VALUE"},
			wantBody: [][]string{
				{"Bucket name", "happy-bucket"},
				{"Region", "ap-northeast-1"},
				{"Block Public Access", "Enabled"},
				{"Encryption", "AES256"},
				{"Versioning", "Enabled"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &initS3Result{
				BucketName:        tt.fields.BucketName,
				Region:            tt.fields.Region,
				BlockPublicAccess: tt.fields.BlockPublicAccess,
				Encryption:        tt.fields.Encryption,
				Versioning:        tt.fields.Versioning,
			}
			gotHeader, gotBody := i.createTableInput()
			if !reflect.DeepEqual(gotHeader, tt.wantHeader) {
				t.Errorf("initS3Result.createTableInput() gotHeader = %v, want %v", gotHeader, tt.wantHeader)
			}
			if !reflect.DeepEqual(gotBody, tt.wantBody) {
				t.Errorf("initS3Result.createTableInput() gotBody = %v, want %v", gotBody, tt.wantBody)
			}
		})
	}
}
