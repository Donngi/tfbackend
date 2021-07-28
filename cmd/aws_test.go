package cmd

import (
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
