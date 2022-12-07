package test

import (
	"testing"
	"web.misaki.world/FinalExam/aws"
)

func TestUploadFile(t *testing.T) {
	type args struct {
		filePath string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "case 1",
			args: args{
				filePath: "20403413/audio/1.mp4",
			},
			want:    "s3://audio-2022-12-08/20403413/audio/1.mp4",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := aws.UploadFile(tt.args.filePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("UploadFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UploadFile() got = %v, want %v", got, tt.want)
			}
		})
	}
}
