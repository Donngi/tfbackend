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
