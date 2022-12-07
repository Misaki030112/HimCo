package test

import (
	"log"
	"testing"
	"web.misaki.world/FinalExam/aws"
)

func TestConvertToText(t *testing.T) {
	log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)
	type args struct {
		mp4Url   string
		savePath string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{
			name: "case 1",
			args: args{
				mp4Url:   "s3://audio-2022-12-08/20403413/audio/1.mp4",
				savePath: "20403413/data",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			aws.ConvertToText(tt.args.mp4Url, tt.args.savePath)
		})
	}
}
